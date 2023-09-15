package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"

	"github.com/lruggieri/fxnow/common/cache/redis"
	cHttp "github.com/lruggieri/fxnow/common/http"
	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/logger/zap"
	"github.com/lruggieri/fxnow/common/store"
	"github.com/lruggieri/fxnow/common/store/mysql"
	"github.com/lruggieri/fxnow/fxrate/logic"
)

var (
	l logic.Logic

	str store.Store
)

func main() {
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

	cache := &redis.Cacher{
		Client: redis.NewClient(redis.Config{
			Addrs: []string{os.Getenv("REDIS_ADDRS")},
		}),
	}

	l = &logic.Impl{
		Store: str,
		Cache: cache,
		Clock: clock.New(),
	}

	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/health", HandleHealth)

	r.GET("/rate", HandleGetRate)

	panic(r.Run(fmt.Sprintf(":%s", port)))
}

func HandleHealth(c *gin.Context) {
	cHttp.HTTPResponse(c, "OK", nil, http.StatusOK)
}

func HandleGetRate(c *gin.Context) {
	now := time.Now()

	from := c.Query("from")
	to := c.Query("to")
	apiKey := c.Query("api-key")

	if from == "" {
		cHttp.HTTPResponse(c, nil, fmt.Errorf("missing 'from' parameter"), http.StatusBadRequest)

		return
	}

	if to == "" {
		cHttp.HTTPResponse(c, nil, fmt.Errorf("missing 'to' parameter"), http.StatusBadRequest)

		return
	}

	if apiKey == "" {
		cHttp.HTTPResponse(c, nil, fmt.Errorf("missing 'api-key' parameter"), http.StatusBadRequest)

		return
	}

	res, err := l.GetRate(context.WithValue(c, logic.ContextKeyAPIKey, apiKey), logic.GetRateRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})
	if err != nil {
		cHttp.HTTPResponse(c, "", err, cHttp.GetHttpStatusFromError(err))

		return
	}

	cHttp.HTTPResponse(c, struct {
		From      string  `json:"from"`
		To        string  `json:"to"`
		Rate      float64 `json:"rate"`
		Timestamp int64   `json:"timestamp"`
		Took      int64   `json:"took"`
	}{
		From:      res.FromCurrency,
		To:        res.ToCurrency,
		Rate:      res.Rate,
		Timestamp: res.Timestamp,
		Took:      time.Since(now).Milliseconds(),
	}, nil, http.StatusOK)
}
