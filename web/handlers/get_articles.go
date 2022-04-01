package handlers

import (
	"github.com/gin-gonic/gin"
	"main/database"
	dbStructure "main/database/structure"
	"main/web"
	"net/http"
	"strconv"
)

func GetArticles(ctx *gin.Context) {

	db := database.GetDB()
	var articles []dbStructure.ArticleBriefInfo

	paramPageNo := ctx.DefaultQuery("page", "1")
	paramPageSize := ctx.DefaultQuery("page_size", "20")
	PageNo, err := strconv.Atoi(paramPageNo)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	pageSize, err := strconv.Atoi(paramPageSize)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	db.Table(dbStructure.TableArticles).Offset((PageNo - 1) * pageSize).Limit(pageSize).Find(&articles)
	web.Write(ctx, http.StatusAccepted, articles)

}
