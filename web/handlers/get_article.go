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
		web.WriteBadRequestError(ctx, "Specify article ID")
		return
	}

	article, found := application.GetApplication().GetArticlesRepo().GetArticle(articleId)
	if !found {
		web.WriteMessage(ctx, http.StatusNotFound, "Article not found")
		return
	}
	web.Write(ctx, http.StatusOK, article)
}
