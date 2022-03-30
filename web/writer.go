package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/models"
	"net/http"
)

func writeNewLine(ctx *gin.Context) {
	_, err := ctx.Writer.Write([]byte("\r\n"))
	if err != nil {
		log.Println(err)
	}
}

func WriteBadRequestError(ctx *gin.Context, message string, prevError ...error) {
	ctx.IndentedJSON(http.StatusBadRequest, &models.Response{
		Code: http.StatusBadRequest,
		Body: message,
	})
	log.Println(message)
	for e := range prevError {
		log.Println(e)
	}
	writeNewLine(ctx)
}

func WriteSuccessMessage(ctx *gin.Context, message string) {
	ctx.IndentedJSON(http.StatusOK, &models.Response{
		Code: http.StatusOK,
		Body: message,
	})
	log.Println(message)
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
	log.Println(object)
	writeNewLine(ctx)
}
