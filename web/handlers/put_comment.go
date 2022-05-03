package handlers

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"on951/application"
	dbStructure "on951/database/structure"
	"on951/image_links_parser"
	"on951/models"
	"on951/strings_manipulator"
	"on951/web"
	"strconv"
	"strings"
)

func PutComment(ctx *gin.Context) {

	authorizedUser, err := application.GetApplication().GetAuthorizedUserFromHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		web.WriteMessage(ctx, http.StatusForbidden, "No logged user!")
		return
	}

	var comment dbStructure.Comment
	articleId, err := strconv.Atoi(strings.TrimSpace(ctx.Param("article_id")))

	if err != nil {
		web.WriteBadRequestError(ctx, "Article ID is missing", err)
		return
	}

	commentBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		web.WriteBadRequestError(ctx, "Error reading comment", err)
		return
	}

	_, found := application.GetApplication().GetArticlesRepo().GetArticle(articleId)

	if !found {
		web.WriteMessage(ctx, http.StatusNotFound, "Article not found!")
		return
	}

	commentBodyStr := string(commentBody)
	var areImageLinksProcessed bool
	commentBodyStr, areImageLinksProcessed = image_links_parser.Process(commentBodyStr, application.GetApplication().GetImagesDir(), application.ImagesDirectory)
	if !areImageLinksProcessed {
		web.Write(ctx, http.StatusInsufficientStorage, models.Response{
			Code: http.StatusInsufficientStorage,
			Body: "Unable to process images in the comment",
		})
		return
	}

	comment.UserId = authorizedUser.Id
	comment.ArticleId = uint(articleId)
	comment.Content = strings_manipulator.StripTags(commentBodyStr)
	created := application.GetApplication().GetArticlesRepo().PutComment(&comment)
	if !created {
		web.WriteMessage(ctx, http.StatusTooEarly, "failed to insert your comment")
		return
	}
	ctx.Header("Content-Location", "/comment/"+strconv.Itoa(int(comment.Id)))
	web.WriteSuccessfullyCreatedMessage(ctx, "Thanks for commenting!")
}
