package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"on951/models"
	"on951/web"
)

const MsgAllGood = "All good"

func HealthCheck(ctx *gin.Context) {
	var rsp = models.Response{
		Code: 200,
		Body: MsgAllGood,
	}
	web.Write(ctx, http.StatusOK, rsp)
}
