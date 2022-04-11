package handlers

import (
	"github.com/gin-gonic/gin"
	"main/state"
	"main/web"
	"net/http"
	"strconv"
	"strings"
)

func GetArticle(ctx *gin.Context) {

	paramArticleId := strings.TrimSpace(ctx.Query("article_id"))
	if paramArticleId == "" {
		web.WriteBadRequestError(ctx, "Specify article ID")
		return
	}
	articleId, err := strconv.Atoi(paramArticleId)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	article, found := state.GetApplication().GetArticlesRepo().GetArticle(articleId)
	if !found {
		web.Write(ctx, http.StatusOK, []string{})
		return
	}
	web.Write(ctx, http.StatusOK, article)
}
