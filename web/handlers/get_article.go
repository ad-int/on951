package handlers

import (
	"github.com/gin-gonic/gin"
	"main/database"
	"main/models"
	"main/web"
	"net/http"
	"strconv"
)

func GetArticle(ctx *gin.Context) {

	db := database.GetDB()
	var article models.Article

	if ctx.Query("id") == "" {
		web.WriteError(ctx, http.StatusBadRequest, "Specify article ID")
		return
	}

	articleId, err := strconv.Atoi(ctx.Query("id"))
	if articleId < 1 || err != nil {
		web.WriteError(ctx, http.StatusBadRequest, "Invalid article ID")
		return
	}

	tx := db.First(&article, articleId)
	if tx.RowsAffected != 1 {
		web.Write(ctx, http.StatusOK, []string{})
		return
	}
	web.Write(ctx, http.StatusOK, article)
}
