package api

import (
	"main/database"
	dbStructure "main/database/structure"
)

type ArticlesRepository interface {
	GetArticles(pageNo int, pageSize int) []dbStructure.ArticleBriefInfo
	GetArticle(articleId int) (dbStructure.ArticleWithComments, bool)
	PutComment(articleId int, comment string) bool
}
type TArticlesRepository struct {
	database.TDatabase
}

func (aRepo *TArticlesRepository) GetArticles(pageNo int, pageSize int) []dbStructure.ArticleBriefInfo {
	var list []dbStructure.ArticleBriefInfo
	aRepo.GetDB().Table(dbStructure.TableArticles).Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&list)
	return list
}

func (aRepo *TArticlesRepository) GetArticle(articleId int) (dbStructure.Article, bool) {
	var article dbStructure.Article
	tx := aRepo.GetDB().First(&article, articleId)
	return article, tx.RowsAffected > 1
}
func (aRepo *TArticlesRepository) GetArticleWithComments(articleId int) (dbStructure.ArticleWithComments, bool) {
	var article dbStructure.ArticleWithComments
	tx := aRepo.GetDB().Debug().
		Table(dbStructure.TableArticles).
		Preload("Comments").
		First(&article, articleId)
	return article, tx.RowsAffected > 1
}
