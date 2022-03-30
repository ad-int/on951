package handlers

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"main/database"
	"main/image_links_parser"
	"main/models"
	"main/web"
	"net/http"
	"strconv"
	"strings"
)

func PutComment(ctx *gin.Context) {

	db := database.GetDB()
	var comment models.Comment
	articleId, err := strconv.Atoi(strings.TrimSpace(ctx.Query("article_id")))

	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	commentBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		web.WriteBadRequestError(ctx, "Error reading comment", err)
		return
	}
	commentBodyStr := string(commentBody)
	var areLmageLinksProcessed bool
	commentBodyStr, areLmageLinksProcessed = image_links_parser.Process(commentBodyStr)
	if !areLmageLinksProcessed {
		web.Write(ctx, http.StatusInsufficientStorage, models.Response{
			Code: http.StatusInsufficientStorage,
			Body: "Unable to process images in the comment",
		})
		return
	}

	comment.ArticleId = uint(articleId)
	comment.Content = commentBodyStr
	tx := db.Create(&comment)
	if tx.Error != nil {
		web.Write(ctx, http.StatusTooEarly, err)
		return
	}
	if tx.RowsAffected < 1 {
		ctx.Status(http.StatusTooEarly)
		return
	}
	ctx.Header("Content-Location", "/comment/"+strconv.Itoa(int(comment.Id)))
	web.WriteSuccessfullyCreatedMessage(ctx, "Thanks for commenting!")
}
