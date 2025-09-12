package server

import (
	"github.com/gin-gonic/gin"
)

// Response object as HTTP response
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// HTTPRes normalize HTTP Response format
func HTTPRes(c *gin.Context, httpCode int, msg string, data any) {
	c.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  msg,
		Data: data,
	})
}
