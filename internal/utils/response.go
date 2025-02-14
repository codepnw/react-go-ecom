package utils

import (
	"github.com/gin-gonic/gin"
)

type IResponse interface {
	Success(code int, data any)
	Error(code int, err error)
}

type Response struct {
	Context *gin.Context
}

func NewResponse(c *gin.Context) IResponse {
	return &Response{Context: c}
}

func (r *Response) Success(code int, data any) {
	switch v := data.(type) {
	case string:
		r.Context.JSON(code, gin.H{"message": v})
	default:
		r.Context.JSON(code, gin.H{"data": v})
	}
}

func (r *Response) Error(code int, err error) {
	r.Context.JSON(code, gin.H{"error": err.Error()})
}
