package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"

	"github.com/lruggieri/fxnow/common/cache/redis"
	cHttp "github.com/lruggieri/fxnow/common/http"
	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/logger/zap"
	"github.com/lruggieri/fxnow/common/store"
	"github.com/lruggieri/fxnow/common/store/mysql"
	"github.com/lruggieri/fxnow/common/util"

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
	r.GET("/fxrate/health", HandleHealth)

	v1 := r.Group("/fxrate/v1")
	v1.GET("/rate", HandleGetRate)

	panic(r.Run(fmt.Sprintf(":%s", port)))
}

func HandleHealth(c *gin.Context) {
	cHttp.HTTPResponse(c, "OK", nil, http.StatusOK)
}

func HandleGetRate(c *gin.Context) {
	now := time.Now()

	// pairs: a list of currency pairs separated by comma (e.g. "USD_JPY,EUR_USD,GBP_CAD")
	pairsStr := c.Query("pairs")

	apiKey := c.Query("api-key")

	if pairsStr == "" {
		cHttp.HTTPResponse(c, nil, fmt.Errorf("missing 'pairs' parameter"), http.StatusBadRequest)

		return
	}

	if apiKey == "" {
		cHttp.HTTPResponse(c, nil, fmt.Errorf("missing 'api-key' parameter"), http.StatusBadRequest)

		return
	}

	pairs := strings.Split(pairsStr, ",")

	cleanPairs := util.Map(pairs, func(item string) string {
		return strings.TrimSpace(item)
	})

	cleanPairs = util.Filter(cleanPairs, func(item string) bool {
		return item != ""
	})

	res, err := l.GetRate(context.WithValue(c, logic.ContextKeyAPIKey, apiKey), logic.GetRateRequest{
		Pairs: cleanPairs,
	})
	if err != nil {
		cHttp.HTTPResponse(c, "", err, cHttp.GetHttpStatusFromError(err))

		return
	}

	type responseRate struct {
		Pair      string  `json:"pair"`
		Rate      float64 `json:"rate"`
		Timestamp int64   `json:"timestamp"`
	}

	rates := make([]responseRate, 0, len(res.Rates))

	for _, rate := range res.Rates {
		rates = append(rates, responseRate{
			Pair:      rate.Pair,
			Rate:      rate.Rate,
			Timestamp: rate.Timestamp,
		})
	}

	cHttp.HTTPResponse(c, struct {
		Rates []responseRate `json:"rates"`
		Took  int64          `json:"took"`
	}{
		Rates: rates,
		Took:  time.Since(now).Milliseconds(),
	}, nil, http.StatusOK)
}
