package data_generator

import (
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"on951/api"
	"on951/application"
	"on951/database"
	dbStructure "on951/database/structure"
	"testing"
)

type dataGeneratorTestSuite struct {
	suite.Suite
	db           database.TDatabaseMock
	recorder     *httptest.ResponseRecorder
	articlesRepo api.ArticlesRepository
}

func TestDataGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(dataGeneratorTestSuite))
}

func (suite *dataGeneratorTestSuite) SetupSuite() {
	suite.db.On("ConnectToDB", "dummy").Return(true)
	suite.Assert().True(suite.db.ConnectToDB("dummy"), "Connecting to in-memory DB")
	suite.db.AutoMigrate()
	suite.db.GetDB().AllowGlobalUpdate = true
	suite.articlesRepo = &api.TArticlesRepository{IDatabase: &suite.db}
	app := &application.TApplicationMock{}
	app.SetDB(&suite.db)
	app.SetArticlesRepo(suite.articlesRepo)
	application.SetApplication(app)
}

func (suite *dataGeneratorTestSuite) TearDownSuite() {
	suite.db.DisconnectDB()
	suite.Nil(application.GetApplication().GetDatabase().GetDB(), "disconnecting in-memory DB")
}

func (suite *dataGeneratorTestSuite) TestGenerateArticle() {
	suite.recorder = httptest.NewRecorder()

	suite.IsType(&application.TApplicationMock{}, application.GetApplication())
	for i := 0; i < 17; i++ {
		GenerateArticle()
	}
	articles := application.GetApplication().GetArticlesRepo().GetArticles(1, 20)
	suite.Len(articles, 17)
}

func (suite *dataGeneratorTestSuite) TestGenerateUser() {
	err := GenerateUser("test-user", "test-password1$", 14)
	suite.Nil(err)
}

func (suite *dataGeneratorTestSuite) TestGenerateUserWithCostTooHigh() {
	err := GenerateUser("test-userH", "test-password1$", 9999)
	suite.NotNil(err)
}

func (suite *dataGeneratorTestSuite) TestGenerateUserWithEmptyName() {
	_ = suite.db.GetDB().Migrator().DropTable(dbStructure.User{})
	_ = suite.db.GetDB().Migrator().AutoMigrate()
	err := GenerateUser("", "test-password1$", 14)
	suite.NotNil(err)
}
