package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/database"
	dbStructure "main/database/structure"
	"main/web"
	"net/http"
	"strconv"
	"strings"
)

func GetComments(ctx *gin.Context) {

	db := database.GetDB()
	var article dbStructure.ArticleWithComments
	articleId := strings.TrimSpace(ctx.Query("article_id"))
	if articleId == "" {
		web.WriteBadRequestError(ctx, "Specify article ID")
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
	tx := db.Debug().
		Table(dbStructure.TableArticles).
		Preload("Comments").
		First(&article, articleId)
	if tx.RowsAffected < 1 {
		web.WriteMessage(ctx, http.StatusNotFound, "Empty :(")
		return
	}
	web.Write(ctx, http.StatusAccepted, article)

}
