package handlers

import (
	"github.com/gin-gonic/gin"
	"main/models"
	"main/web"
	"net/http"
)

const MsgAllGood = "All good"

func HealthCheck(ctx *gin.Context) {
	var rsp = models.Response{
		Code: 202,
		Body: MsgAllGood,
	}
	web.Write(ctx, http.StatusAccepted, rsp)
}
