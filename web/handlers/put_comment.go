package handlers

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"main/database"
	dbStructure "main/database/structure"
	"main/image_links_parser"
	"main/models"
	"main/state"
	"main/web"
	"net/http"
	"strconv"
	"strings"
)

func PutComment(ctx *gin.Context) {

	authorizedUser, err := state.GetAuthorizedUserFromHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		web.WriteMessage(ctx, http.StatusForbidden, "No logged user!")
		return
	}

	db := database.GetDB()
	var comment dbStructure.Comment
	articleId, err := strconv.Atoi(strings.TrimSpace(ctx.Query("article_id")))

	if err != nil {
		web.WriteBadRequestError(ctx, "Article ID is missing", err)
		return
	}

	commentBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		web.WriteBadRequestError(ctx, "Error reading comment", err)
		return
	}
	commentBodyStr := string(commentBody)
	var areImageLinksProcessed bool
	commentBodyStr, areImageLinksProcessed = image_links_parser.Process(commentBodyStr)
	if !areImageLinksProcessed {
		web.Write(ctx, http.StatusInsufficientStorage, models.Response{
			Code: http.StatusInsufficientStorage,
			Body: "Unable to process images in the comment",
		})
		return
	}

	comment.UserId = authorizedUser.Id
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
