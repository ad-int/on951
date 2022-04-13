package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"on951/application"
	dbStructure "on951/database/structure"
	"on951/web"
	"strconv"
)

func GetArticles(ctx *gin.Context) {
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
	var articles []dbStructure.ArticleBriefInfo
	articles = application.GetApplication().GetArticlesRepo().GetArticles((PageNo-1)*pageSize, pageSize)
	web.Write(ctx, http.StatusAccepted, articles)

}
