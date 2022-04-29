package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"on951/application"
	"on951/web"
	"strconv"
	"strings"
)

func GetComments(ctx *gin.Context) {
	var err error
	pageNo := defaultPageNo
	pageSize := defaultPageSize

	paramArticleId := strings.TrimSpace(ctx.Param("article_id"))
	if paramArticleId == "" {
		web.WriteBadRequestError(ctx, "Specify article ID")
		return
	}
	articleId, err := strconv.Atoi(paramArticleId)
	if err != nil {
		web.WriteBadRequestError(ctx, "Incorrect article ID", err)
		return
	}
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
	article, found := application.GetApplication().GetArticlesRepo().GetArticleWithComments(articleId, pageNo, pageSize)
	if !found {
		web.WriteMessage(ctx, http.StatusNotFound, "Empty :(")
		return
	}
	web.Write(ctx, http.StatusOK, article)

}
