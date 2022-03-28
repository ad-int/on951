package web

import (
	"github.com/gin-gonic/gin"
	"main/models"
)

func WriteError(ctx *gin.Context, code int, message string) {
	ctx.IndentedJSON(code, &models.Error{
		Code:    code,
		Message: message,
	})
	ctx.Writer.Write([]byte("\r\n"))
}

func Write(ctx *gin.Context, code int, object interface{}) {
	ctx.IndentedJSON(code, object)
	ctx.Writer.Write([]byte("\r\n"))
}
