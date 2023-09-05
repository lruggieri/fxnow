package auth

import "context"

const (
	ContextUserInfoKey = "user-info"
)

type Authenticator interface {
	GetOIDCConsentURL(redirectURL string) string
	AuthenticateOIDC(code string) (string, error)

	IsJWTValid(token string) bool
	GetUserInfo(token string) *UserInfo
}

type Config struct {
	OIDC OIDCConfig
}

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type UserInfo struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

func GetUserInfoFromContext(ctx context.Context) *UserInfo {
	ctxInfo := ctx.Value(ContextUserInfoKey)

	userInfo, ok := ctxInfo.(*UserInfo)
	if !ok {
		return nil
	}

	return userInfo
}
