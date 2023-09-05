package http

import "github.com/gin-gonic/gin"

type Response struct {
	Response interface{} `json:"response,omitempty"`
	Error    string      `json:"error,omitempty"`
}

func HTTPResponse(c *gin.Context, resp interface{}, err error, statusCode int) {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	c.JSON(statusCode, Response{Response: resp, Error: errStr})
}
