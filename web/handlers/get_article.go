package handlers

import (
	"github.com/gin-gonic/gin"
	"main/database"
	dbStructure "main/database/structure"
	"main/web"
	"net/http"
	"strings"
)

func GetArticle(ctx *gin.Context) {

	db := database.GetDB()
	var article dbStructure.Article

	articleId := strings.TrimSpace(ctx.Query("article_id"))
	if articleId == "" {
		web.WriteBadRequestError(ctx, "Specify article ID")
		return
	}

	tx := db.First(&article, articleId)
	if tx.RowsAffected != 1 {
		web.Write(ctx, http.StatusOK, []string{})
		return
	}
	web.Write(ctx, http.StatusOK, article)
}
