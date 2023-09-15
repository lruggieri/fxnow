package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/lruggieri/fxnow/common/cache/redis"
	"github.com/lruggieri/fxnow/common/client/fastforex"
	cHttp "github.com/lruggieri/fxnow/common/http"
	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/logger/zap"

	"github.com/lruggieri/fxnow/fxrate/logic"
)

var l logic.Logic

func main() {
	mainContext := context.Background()

	logger.InitLogger(zap.New(zap.Config{
		Development: false,
		Level:       logger.LevelDebug,
	}))

	port := os.Getenv("PORT")

	cache := &redis.Cacher{
		Client: redis.NewClient(redis.Config{
			Addrs: []string{os.Getenv("REDIS_ADDRS")},
		}),
	}

	l = &logic.Impl{
		Cache: cache,
		FXSource: fastforex.NewClient(
			os.Getenv("FASTFOREX_API_KEY"),
			http.DefaultClient,
		),
	}

	// start service logic
	go l.StartFXUpdate(mainContext)

	r := gin.Default()
	r.GET("/health", HandleHealth)

	panic(r.Run(fmt.Sprintf("127.0.0.1:%s", port)))
}

func HandleHealth(c *gin.Context) {
	cHttp.HTTPResponse(c, "OK", nil, http.StatusOK)
}
