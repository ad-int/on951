package handlers

import (
	"github.com/gin-gonic/gin"
	"on951/database/structure"
	"on951/models"
)

type Handlers = struct {
	GetArticle  []TestCase
	GetArticles []TestCase
	GetComments []TestCase
	PutComment  []TestCase
}
type TestCase struct {
	requestURI        string
	totalArticlesInDb int
	params            gin.Params
	response          interface{}
}

var testHandlersData = Handlers{
	GetArticle: []TestCase{
		{
			requestURI:        "",
			totalArticlesInDb: 20,
			params: gin.Params{
				{Key: "article_id", Value: "3"},
			},
			response: structure.Article{
				Id:          3,
				AuthorId:    0,
				Title:       "article 3",
				Description: "",
			},
		},
		{
			requestURI: "",
			params:     gin.Params{},
			response: models.Response{
				Code: 400,
				Body: "Specify article ID",
			},
		},
	},
	GetArticles: []TestCase{
		{
			requestURI:        "",
			params:            gin.Params{},
			totalArticlesInDb: 5,
			response: []structure.ArticleBriefInfo{
				{
					Id:    1,
					Title: "article 1",
				},
				{
					Id:    2,
					Title: "article 2",
				},
				{
					Id:    3,
					Title: "article 3",
				},
				{
					Id:    4,
					Title: "article 4",
				},
				{
					Id:    5,
					Title: "article 5",
				},
			},
		},
	},
	GetComments: []TestCase{
		{
			requestURI: "",
			params: gin.Params{
				{Key: "article_id", Value: "3"},
			},
			totalArticlesInDb: 5,
			response: structure.ArticleWithComments{

				Id:          3,
				Title:       "article 3",
				AuthorId:    0,
				Description: "",
				Comments:    []structure.Comment{},
			},
		},
	},
	PutComment: []TestCase{},
}
