package handlers

import (
	"github.com/gin-gonic/gin"
	dbStructure "main/database/structure"
	"main/state"
	"main/web"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultParamPageNo   = "1"
	defaultParamPageSize = "20"
)

func GetArticles(ctx *gin.Context) {
	paramPageNo := strings.TrimSpace(ctx.Param("page"))
	paramPageSize := strings.TrimSpace(ctx.Param("page_size"))

	if len(paramPageNo) < 1 {
		paramPageNo = defaultParamPageNo

	}
	if len(paramPageSize) < 1 {
		paramPageNo = defaultParamPageSize

	}

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
