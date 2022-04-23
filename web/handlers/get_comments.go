package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"on951/application"
	"on951/web"
	"strconv"
	"strings"
)

func GetComments(ctx *gin.Context) {
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
	paramPageNo := ctx.DefaultQuery("page", "1")
	paramPageSize := ctx.DefaultQuery("page_size", "20")
	PageNo, err := strconv.Atoi(paramPageNo)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	pageSize, err := strconv.Atoi(paramPageSize)
	log.Println(pageSize, PageNo)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err != nil {
		web.Write(ctx, http.StatusOK, err)
		return
	}
	article, found := application.GetApplication().GetArticlesRepo().GetArticleWithComments(articleId)
	if !found {
		web.WriteMessage(ctx, http.StatusNotFound, "Empty :(")
		return
	}
	web.Write(ctx, http.StatusAccepted, article)

}
