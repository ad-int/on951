package api

import dbStructure "on951/database/structure"

type ArticlesRepository interface {
	GetArticles(pageNo int, pageSize int) []dbStructure.ArticleBriefInfo
	GetArticle(articleId int) (dbStructure.Article, bool)
	GetArticleWithComments(articleId int, pageNo int, pageSize int) (dbStructure.ArticleWithComments, bool)
	PutComment(comment *dbStructure.Comment) bool
}
