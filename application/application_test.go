package application

import (
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"log"
	"on951/api"
	"on951/database"
	"os"
	"reflect"
	"testing"
)

const fakeImagesDir = "path/fake-images-dir"

type applicationTestSuite struct {
	suite.Suite
	config       map[string]string
	articlesRepo *api.TArticlesRepository
}

func TestApplicationTestSuite(t *testing.T) {
	suite.Run(t, new(applicationTestSuite))
}

func (suite *applicationTestSuite) TestGetApplication() {
	suite.IsType(&TApplication{}, GetApplication())
}

func (suite *applicationTestSuite) TestSetApplication() {

	suite.IsType(&TApplication{}, GetApplicationRepository().GetApplication())

	SetApplication(&TApplicationMock{})
	suite.IsType(&TApplicationMock{}, GetApplicationRepository().GetApplication())
}

func (suite *applicationTestSuite) TestSetArticlesRepo() {

	app := &TApplication{}
	app.SetArticlesRepo(suite.articlesRepo)

	suite.Same(suite.articlesRepo, app.GetArticlesRepo())
}

func (suite *applicationTestSuite) TestReadEnvFile() {

	var read bool
	var config map[string]string

	app := &TApplication{ConfigFilePath: "fdfedgdf"}

	suite.Panics(func() {
		_, _ = app.ReadEnvFile()
	})
	suite.False(read)
	f, err := ioutil.TempFile("", "config-test")
	suite.Nil(err)
	_, _ = f.WriteString(`TEST1="value1"` + "\n" + `TEST2="value2"`)
	app = &TApplication{ConfigFilePath: f.Name()}
	read, config = app.ReadEnvFile()
	_ = f.Close()
	suite.True(read)
	suite.Equal(map[string]string{"TEST1": "value1", "TEST2": "value2"}, config)
}

func (suite *applicationTestSuite) TestInitDb() {
	dbMock := &database.TDatabaseMock{}
	dbMock.On("GetConfigValue", "DSN").Return("")
	dbMock.On("ConnectToDB", "").Return(true)
	app := &TApplication{db: dbMock}
	suite.True(app.InitDb())
}

func (suite *applicationTestSuite) TestInit() {
	dbMock := &database.TDatabaseMock{}
	dbMock.On("GetConfigValue", "DSN").Return("")
	dbMock.On("ConnectToDB", "").Return(true)
	app := &TApplicationMock{db: dbMock}
	app.On("Init", &testApplicationData[0].routes).Panic(MsgCannotReadImagesDirectory)
	suite.PanicsWithValue(MsgCannotReadImagesDirectory, func() {
		_ = app.Init(&testApplicationData[0].routes)
	})
}

func (suite *applicationTestSuite) TestInitDbFailsWithPanic() {
	dbMock := &database.TDatabaseMock{}
	dbMock.On("GetConfigValue", "DSN").Return("invalid-dsn")
	dbMock.On("ConnectToDB", "invalid-dsn").Return(false)

	app := &TApplicationMock{db: dbMock}
	app.On("GetConfigValue", "DSN").Return("invalid-dsn")
	app.On("InitDb").Panic("could not connect to DB")
	suite.PanicsWithValue("could not connect to DB", func() { app.InitDb() })
}

func (suite *applicationTestSuite) TestApplicationBootstrap() {

	for _, testCase := range testApplicationData {
		dbMock := &database.TDatabaseMock{}
		dbMock.On("GetConfigValue", "DSN").Return("")
		dbMock.On("ConnectToDB", "").Return(true)
		var err error
		imagesTmpDir := testCase.imagesDir
		if len(testCase.imagesDir) > 0 {
			imagesTmpDir, err = ioutil.TempDir(os.TempDir(), testCase.imagesDir)
			suite.Nil(err, "Creating images temp dir")
			log.Println(imagesTmpDir)
		}
		app := &TApplicationMock{db: dbMock, ImagesDir: imagesTmpDir}

		app.On("ReadEnvFile").Return(len(app.config) > 0, testCase.config)
		app.On("GetConfigValue", "DSN").Return("")
		app.On("GetConfigValue", "TRUSTED_PROXIES").Return(testCase.config["TRUSTED_PROXIES"])
		app.On("GetConfigValue", "CORS_ALLOWED_HEADERS").Return(testCase.config["CORS_ALLOWED_HEADERS"])
		app.On("GetConfigValue", "CORS_ALLOW_ALL_ORIGINS").Return(testCase.config["CORS_ALLOW_ALL_ORIGINS"])
		app.On("GetImagesDir").Return(testCase.imagesDir)
		app.On("Init", &testCase.routes).Return(nil)
		app.On("InitDb").Return(true)
		app.On("GetRouter").Return(&app.router)
		app.On("Init").Return(nil)
		appRepo := &TApplicationRepository{application: app}
		log.Println(reflect.TypeOf(appRepo.GetApplication().GetRouter()))

		suite.IsType(&TApplicationMock{}, appRepo.GetApplication())
		suite.Equal(testCase.imagesDir, appRepo.GetApplication().GetImagesDir())

		testFuncInsideBootstrap := ""
		if len(testCase.imagesDir) == 0 {
			suite.PanicsWithError(MsgCannotReadImagesDirectory, func() {
				appRepo.Bootstrap(&testCase.routes)

			})
			continue
		} else if testCase.configurationFails {
			suite.Panics(func() {
				appRepo.Bootstrap(&testCase.routes)
			})
			continue
		} else {
			appRepo.Bootstrap(&testCase.routes, func() { testFuncInsideBootstrap = "test" })

		}

		routesFound := 0
		suite.Equal("test", testFuncInsideBootstrap)
		for method, routeDescription := range testCase.routes {
			for path, _ := range *routeDescription {
				for _, h := range appRepo.GetApplication().GetRouter().GetEngine().Routes() {
					if h.Method == method && h.Path == "/"+path {
						routesFound++
						suite.T().Logf("Added route: %v %v\n", method, path)
					}
				}
			}
		}
		suite.Equal(len(appRepo.GetApplication().GetRouter().GetEngine().Routes()), routesFound)
	}
}
