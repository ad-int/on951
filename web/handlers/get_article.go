package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"on951/application"
	"on951/web"
	"strconv"
	"strings"
)

func GetArticle(ctx *gin.Context) {

	paramArticleId := strings.TrimSpace(ctx.Param("article_id"))
	if paramArticleId == "" {
		web.WriteBadRequestError(ctx, "Specify article ID")
		return
	}
	articleId, err := strconv.Atoi(paramArticleId)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	article, found := application.GetApplication().GetArticlesRepo().GetArticle(articleId)
	if !found {
		web.Write(ctx, http.StatusOK, []string{})
		return
	}
	web.Write(ctx, http.StatusOK, article)
}
