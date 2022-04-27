package api

import (
	"github.com/stretchr/testify/mock"
	"math/rand"
	"on951/database"
	dbStructure "on951/database/structure"
	"on951/randomizer"
	"time"
)

type TArticlesRepositoryMock struct {
	mock.Mock
	database.IDatabase
	TArticlesRepository
}

func (aRepo *TArticlesRepositoryMock) GetArticle(articleId int) (dbStructure.Article, bool) {
	rand.Seed(time.Now().UnixNano())
	return dbStructure.Article{
		Title:       randomizer.GetRandomString(16, 1),
		Description: randomizer.GetRandomString(255, 8),
		AuthorId:    uint(rand.Intn(50)) + 1,
	}, true
}

func (aRepo *TArticlesRepositoryMock) PutComment(comment *dbStructure.Comment) bool {
	return false
}
