package handlers

import (
	"github.com/gin-gonic/gin"
	"main/database"
	"main/models"
	"net/http"
	"strconv"
)

func GetArticles(ctx *gin.Context) {

	db := database.GetDB()
	// Read
	var articles []models.Article

	paramPageNo := ctx.DefaultQuery("page", "1")
	paramPageSize := ctx.DefaultQuery("page_size", "20")
	PageNo, err := strconv.Atoi(paramPageNo)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	pageSize, err := strconv.Atoi(paramPageSize)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	db.Offset((PageNo - 1) * pageSize)
	db.Limit(pageSize)
	db.Find(&articles)



	ctx.JSON(http.StatusAccepted, articles)
}
