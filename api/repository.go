package api

import (
	"fmt"
	"on951/database"
	dbStructure "on951/database/structure"
)

type TArticlesRepository struct {
	database.IDatabase
	ArticlesRepository
}

func (aRepo *TArticlesRepository) GetArticles(pageNo int, pageSize int) []dbStructure.ArticleBriefInfo {
	offset := 0
	if pageNo >= 1 {
		offset = (pageNo - 1) * pageSize
	}
	var list []dbStructure.ArticleBriefInfo
	aRepo.GetDB().
		Debug().
		Table(dbStructure.TableArticles).
		Offset(offset).
		Limit(pageSize).
		Find(&list)
	return list
}

func (aRepo *TArticlesRepository) GetArticle(articleId int) (dbStructure.Article, bool) {
	var article dbStructure.Article
	tx := aRepo.GetDB().First(&article, articleId)
	return article, tx.RowsAffected > 0
}
func (aRepo *TArticlesRepository) GetArticleWithComments(articleId int, pageNo int, pageSize int) (dbStructure.ArticleWithComments, bool) {
	var article dbStructure.ArticleWithComments
	offset := 0
	if pageNo >= 1 {
		offset = (pageNo - 1) * pageSize
	}
	tx := aRepo.GetDB().
		Debug().
		Table(dbStructure.TableArticles).
		Preload("Comments", fmt.Sprintf("1 LIMIT %v OFFSET %v", pageSize, offset)).
		First(&article, articleId)
	return article, tx.RowsAffected > 0
}

func (aRepo *TArticlesRepository) PutComment(comment *dbStructure.Comment) bool {
	tx := aRepo.GetDB().Debug().Create(&comment)
	return tx.RowsAffected > 0
}
