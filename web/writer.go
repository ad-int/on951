package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"on951/models"
)

func writeNewLine(ctx *gin.Context) {
	_, err := ctx.Writer.Write([]byte("\r\n"))
	if err != nil {
		log.Println(err)
	}
}

func WriteBadRequestError(ctx *gin.Context, message string, prevError ...error) {
	WriteMessage(ctx, http.StatusBadRequest, message, prevError...)
}
func WriteMessage(ctx *gin.Context, code int, message string, prevError ...error) {
	ctx.IndentedJSON(code, &models.Response{
		Code: code,
		Body: message,
	})
	log.Println(message)
	for _, e := range prevError {
		if e != nil {
			_ = ctx.Error(e)
		}
		log.Println(e)
	}
	writeNewLine(ctx)
}

func WriteSuccessfullyCreatedMessage(ctx *gin.Context, message string) {
	ctx.IndentedJSON(http.StatusCreated, &models.Response{
		Code: http.StatusCreated,
		Body: message,
	})
	log.Println(message)
	writeNewLine(ctx)
}

func Write(ctx *gin.Context, code int, object interface{}) {
	ctx.IndentedJSON(code, object)
	if code < 200 || code > 299 {
		log.Println(object)
	}
	writeNewLine(ctx)
}
