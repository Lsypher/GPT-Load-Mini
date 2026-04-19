package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int      `json:"code"`
	Data any      `json:"data,omitempty"`
	Msg  string   `json:"message,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(200, Response{Code: 0, Data: data})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{Code: code, Msg: message})
}

func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(200, Response{Code: 0, Msg: message})
}
