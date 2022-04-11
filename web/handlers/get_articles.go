package handlers

import (
	"github.com/gin-gonic/gin"
	dbStructure "main/database/structure"
	"main/state"
	"main/web"
	"net/http"
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
	articles = state.GetApplication().GetArticlesRepo().GetArticles((PageNo-1)*pageSize, pageSize)
	web.Write(ctx, http.StatusAccepted, articles)

}
