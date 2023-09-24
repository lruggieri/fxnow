package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	cHttp "github.com/lruggieri/fxnow/common/http"
	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/logger/zap"
	"github.com/lruggieri/fxnow/common/store"
	"github.com/lruggieri/fxnow/common/store/mysql"

	"github.com/lruggieri/fxnow/identity/auth"
	"github.com/lruggieri/fxnow/identity/logic"
)

var (
	authenticator auth.Authenticator

	l logic.Logic

	str store.Store
)

func main() {
	mainContext := context.Background()

	logger.InitLogger(zap.New(zap.Config{
		Development: false,
		Level:       logger.LevelDebug,
	}))

	port := os.Getenv("PORT")

	mysqlPort, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		panic(err)
	}

	str, err = mysql.New(mysql.Config{
		Username: os.Getenv("MYSQL_USERNAME"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     mysqlPort,
		DBName:   os.Getenv("MYSQL_DB_NAME"),
	})
	if err != nil {
		panic(err)
	}

	l = &logic.Impl{
		Store: str,
	}

	authenticator, err = auth.NewBasic(mainContext, auth.Config{
		OIDC: auth.OIDCConfig{
			ClientID:     os.Getenv("OIDC_CLIENT_ID"),
			ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OIDC_REDIRECT_URL"),
		},
	})

	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(commonHeaders())
	r.GET("/identity/health", HandleHealth)

	// oidc
	r.GET("/identity/access", HandleAccess)
	r.GET("/identity/callback", HandleOauthCallback)

	// API key
	v1 := r.Group("/identity/v1")
	v1.GET("/api-keys", HandleListAPIKey)
	v1.POST("/api-key", HandleCreateAPIKey)
	v1.DELETE("/api-key/:key", HandleRevokeAPIKey)

	panic(r.Run(fmt.Sprintf(":%s", port)))
}

func commonHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func HandleHealth(c *gin.Context) {
	cHttp.HTTPResponse(c, "OK", nil, http.StatusOK)
}

func HandleAccess(c *gin.Context) {
	if isAuthenticated(c) {
		cHttp.HTTPResponse(c, "OK", nil, http.StatusOK)
		return
	}

	// invalid token, redirect to OIDC provider
	c.Redirect(http.StatusFound, authenticator.GetOIDCConsentURL(getFullPath(c)))
}

func HandleOauthCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := authenticator.AuthenticateOIDC(code)
	if err != nil {
		cHttp.HTTPResponse(c, "", err, http.StatusInternalServerError)
		return
	}

	c.SetCookie("access_token", token, 0, "/", "", false, false)

	// redirect the user where it came from before the auth flow started
	state := c.Query("state")
	if state != "" {
		c.Redirect(http.StatusFound, state)
		return
	}

	cHttp.HTTPResponse(c, "OK", nil, http.StatusOK)
}

func HandleListAPIKey(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, authenticator.GetOIDCConsentURL(getFullPath(c)))
		return
	}

	resp, err := l.ListAPIKeys(
		context.WithValue(c, auth.ContextUserInfoKey, authenticator.GetUserInfo(getToken(c))),
		logic.ListAPIKeysRequest{},
	)
	if err != nil {
		cHttp.HTTPResponse(c, "", err, cHttp.GetHttpStatusFromError(err))

		return
	}

	type key struct {
		APIKeyID   string `json:"api_key"`
		Expiration int64  `json:"expiration"`
	}

	apiKeys := make([]key, 0, len(resp.APIKeys))

	for _, apiKey := range resp.APIKeys {
		apiKeys = append(apiKeys, key{
			APIKeyID:   apiKey.APIKeyID,
			Expiration: apiKey.Expiration,
		})
	}

	cHttp.HTTPResponse(c, struct {
		APIKeys []key `json:"api-keys"`
	}{apiKeys}, nil, http.StatusOK)
}

func HandleCreateAPIKey(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, authenticator.GetOIDCConsentURL(getFullPath(c)))
		return
	}

	resp, err := l.CreateAPIKey(
		context.WithValue(c, auth.ContextUserInfoKey, authenticator.GetUserInfo(getToken(c))),
		logic.CreateAPIKeyRequest{},
	)
	if err != nil {
		cHttp.HTTPResponse(c, "", err, cHttp.GetHttpStatusFromError(err))

		return
	}

	cHttp.HTTPResponse(c, struct {
		ID string `json:"id"`
	}{
		ID: resp.APIKeyID,
	}, nil, http.StatusOK)
}

func HandleRevokeAPIKey(c *gin.Context) {
	if !isAuthenticated(c) {
		c.Redirect(http.StatusFound, authenticator.GetOIDCConsentURL(getFullPath(c)))
		return
	}

	keyToRevoke := c.Param("key")
	if len(keyToRevoke) == 0 {
		cHttp.HTTPResponse(c, "", fmt.Errorf("invalid key"), http.StatusBadRequest)

		return
	}

	_, err := l.DeleteAPIKey(
		context.WithValue(c, auth.ContextUserInfoKey, authenticator.GetUserInfo(getToken(c))),
		logic.DeleteAPIKeyRequest{APIKeyID: keyToRevoke},
	)
	if err != nil {
		cHttp.HTTPResponse(c, "", err, cHttp.GetHttpStatusFromError(err))

		return
	}

	cHttp.HTTPResponse(c, nil, nil, http.StatusOK)
}

func isAuthenticated(c *gin.Context) bool {
	accessToken := getToken(c)
	if accessToken != "" && authenticator.IsJWTValid(accessToken) {
		return true
	}

	return false
}

func getToken(c *gin.Context) string {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		return ""
	}

	return accessToken
}

func getFullPath(c *gin.Context) string {
	req := c.Request

	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}

	host := req.Host
	path := req.URL.Path
	fullURL := fmt.Sprintf("%s://%s%s", scheme, host, path)

	return fullURL
}
