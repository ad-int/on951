package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"on951/application"
	dbStructure "on951/database/structure"
	"on951/web"
	"strconv"
	"strings"
)

const (
	defaultPageNo   = 1
	defaultPageSize = 20
)

func GetArticles(ctx *gin.Context) {
	var err error
	pageNo := defaultPageNo
	pageSize := defaultPageSize

	paramPageNo := strings.TrimSpace(ctx.Query("page"))
	paramPageSize := strings.TrimSpace(ctx.Query("page_size"))

	if paramPageNo != "" {
		pageNo, err = strconv.Atoi(paramPageNo)
		if err != nil {
			web.WriteBadRequestError(ctx, "Incorrect page number", err)
			return
		}
	}
	if paramPageSize != "" {
		pageSize, err = strconv.Atoi(paramPageSize)
		if err != nil {
			web.WriteBadRequestError(ctx, "Incorrect page size", err)
			return
		}
	}
	var articles []dbStructure.ArticleBriefInfo
	articles = application.GetApplication().GetArticlesRepo().GetArticles(pageNo, pageSize)
	web.Write(ctx, http.StatusAccepted, articles)

}
