package application

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pascaldekloe/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"on951/api"
	"on951/database"
	dbStructure "on951/database/structure"
	"on951/router"
	"os"
	"reflect"
	"testing"
	"time"
)

type FakeRouter struct {
	mock.Mock
	router.TAppRouter
}

func (rf *FakeRouter) Run(addr ...string) error {
	args := rf.Called(addr)
	return args.Error(0)
}

func (rf *FakeRouter) Configure(env string, trustedProxies []string, allowedHeaders []string, allowAllOrigins bool) error {
	args := rf.Called()
	return args.Error(0)
}
func (rf *FakeRouter) InitRoutes(routes *map[string]router.TRoutesList) {

}

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

func (suite *applicationTestSuite) TestSetDB() {

	app := &TApplicationMock{}
	app.SetDB(&database.TDatabaseMock{})

	suite.IsType(&database.TDatabaseMock{}, app.GetDatabase())
}

func (suite *applicationTestSuite) TestGetDatabase() {

	app := &TApplication{db: &database.TDatabase{}}
	suite.IsType(&database.TDatabase{}, app.GetDatabase())
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
	f, err := ioutil.TempFile(os.TempDir(), "config-test")
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

func (suite *applicationTestSuite) TestInitOnApplicationMock() {
	dbMock := &database.TDatabaseMock{}
	dbMock.On("GetConfigValue", "DSN").Return("")
	dbMock.On("ConnectToDB", "").Return(true)
	app := &TApplicationMock{db: dbMock}
	app.On("Init", &testApplicationData[0].routes).Panic(MsgCannotReadImagesDirectory)
	suite.PanicsWithValue(MsgCannotReadImagesDirectory, func() {
		_ = app.Init(&testApplicationData[0].routes)
	})
}

func (suite *applicationTestSuite) TestInit() {

	dbMock := &database.TDatabaseMock{}
	dbMock.On("ConnectToDB", "").Return(true)
	fakeRouter := &FakeRouter{}
	fakeRouter.On("Configure", mock.Anything).Return(nil)
	fakeRouter.On("Run", mock.Anything).Return(nil)
	app := &TApplication{db: dbMock, router: fakeRouter}
	suite.Nil(app.Init(&testApplicationData[0].routes))
	suite.IsType(&FakeRouter{}, app.GetRouter())
}

func (suite *applicationTestSuite) TestGetImagesDir() {

	app := &TApplication{ImagesDir: "application_test.go"}
	suite.Empty(app.GetImagesDir())
}

func (suite *applicationTestSuite) TestGetAuthorizedUserFromHeader() {

	f, err := ioutil.TempFile(os.TempDir(), "config-auth-test")
	suite.Nil(err, "creating temp application config")
	_, _ = f.WriteString(`SECRET="value1"` + "\n" + `ISSUER="value2"` + "\n" + `AUDIENCE="value3"`)
	app := &TApplication{ConfigFilePath: f.Name()}
	SetApplication(app)
	_, _ = app.ReadEnvFile()

	userRecord := dbStructure.User{
		Name:     "guest",
		Password: "not-set",
	}
	token, storedUserRecord, err := suite.generateToken(
		userRecord,
		`SECRET="v1"`+"\n"+`ISSUER="v2"`+"\n"+`AUDIENCE="v3"`,
		time.Now(),
		time.Now(),
		time.Now().Add(24*time.Hour),
	)
	user, err := app.GetAuthorizedUserFromHeader("Bearer " + token)
	suite.Nil(err, "valid token")
	suite.Equal(storedUserRecord, &user)

}

func (suite *applicationTestSuite) generateToken(subject dbStructure.User, configStr string, now time.Time, notBefore time.Time, expires time.Time) (string, *dbStructure.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(subject.Password), 14)
	storedUserRecord := dbStructure.User{
		Id:       1,
		Name:     subject.Name,
		Password: string(hash),
	}
	userRecordJson, _ := json.Marshal(storedUserRecord)
	var claims jwt.Claims
	claims.Subject = string(userRecordJson)
	claims.Issued = jwt.NewNumericTime(now)
	claims.NotBefore = jwt.NewNumericTime(notBefore)
	claims.Expires = jwt.NewNumericTime(expires)
	claims.Issuer = GetApplication().GetConfigValue("ISSUER")
	claims.Audiences = []string{GetApplication().GetConfigValue("AUDIENCE")}

	var jwtBytes []byte
	jwtBytes, err = claims.HMACSign(jwt.HS512, []byte(GetApplication().GetConfigValue("SECRET")))
	return string(jwtBytes), &storedUserRecord, err
}

func (suite *applicationTestSuite) TestGetAuthorizedUserFromHeaderWithInvalidTokens() {

	f, err := ioutil.TempFile(os.TempDir(), "config-auth-test")
	suite.Nil(err, "creating temp application config")
	_, _ = f.WriteString(`SECRET="v1"` + "\n" + `ISSUER="v2"` + "\n" + `AUDIENCE="v3"`)
	app := &TApplication{ConfigFilePath: f.Name()}
	SetApplication(app)
	_, _ = app.ReadEnvFile()

	userRecord := dbStructure.User{
		Name:     "guest",
		Password: "not-set",
	}
	token, _, err := suite.generateToken(
		userRecord,
		`SECRET="v1"`+"\n"+`ISSUER="v2"`+"\n"+`AUDIENCE="v3"`,
		time.Now(),
		time.Now(),
		time.Now().Add(24*time.Hour),
	)

	expiredToken, _, err := suite.generateToken(
		userRecord,
		`SECRET="v1"`+"\n"+`ISSUER="v2"`+"\n"+`AUDIENCE="v3"`,
		time.Now(),
		time.Now(),
		time.Now().Add(-24*time.Hour),
	)

	for _, configStr := range []string{
		`SECRET="v1"` + "\n" + `ISSUER="value2"` + "\n" + `AUDIENCE="value3"`,
		`SECRET="v1"` + "\n" + `AUDIENCE="value3"`,
		`SECRET="v1"` + "\n" + `ISSUER="value3"`,
		`SECRET="v1"` + "\n" + `AUDIENCE="v3"`,
	} {

		f, err = ioutil.TempFile(os.TempDir(), "config-auth-test-*")
		suite.Nil(err, "creating temp application config")
		_, _ = f.WriteString(configStr)
		app = &TApplication{ConfigFilePath: f.Name()}
		SetApplication(app)
		_, _ = app.ReadEnvFile()

		_, err = app.GetAuthorizedUserFromHeader("")
		suite.NotNil(err, "Invalid token")

		_, err = app.GetAuthorizedUserFromHeader("yo")
		suite.NotNil(err, "Invalid token")

		_, err = app.GetAuthorizedUserFromHeader("yo 123")
		suite.NotNil(err, "Invalid token")

		_, err = app.GetAuthorizedUserFromHeader("Bearer 123")
		suite.NotNil(err, "Invalid token")

		_, err = app.GetAuthorizedUserFromHeader("Bearer " + token)
		suite.NotNil(err, "Invalid token")

		_, err = app.GetAuthorizedUserFromHeader("Bearer " + expiredToken)
		suite.NotNil(err, "Invalid token")
	}

}

func (suite *applicationTestSuite) TestInitWithInvalidConfig() {

	dbMock := &database.TDatabaseMock{}
	dbMock.On("GetConfigValue", "DSN").Return("")
	dbMock.On("ConnectToDB", "").Return(true)
	fakeRouter := &FakeRouter{}
	fakeRouter.On("Configure", mock.Anything).Return(errors.New("configuration failed"))
	app := &TApplication{db: dbMock, router: fakeRouter}
	suite.Error(errors.New("configuration failed"), app.Init(&testApplicationData[0].routes))
}

func (suite *applicationTestSuite) TestInitWithNotExistingImagesDir() {

	dbMock := &database.TDatabaseMock{}
	dbMock.On("GetConfigValue", "DSN").Return("")
	dbMock.On("ConnectToDB", "").Return(true)
	app := &TApplication{db: dbMock, router: &FakeRouter{}, ImagesDir: "fdghgfdhgfdhgfesgera"}

	suite.PanicsWithError(MsgCannotReadImagesDirectory, func() {
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
		if testCase.config != nil {
			for key, value := range testCase.config {
				app.On("GetConfigValue", key).Return(value)
			}
		}
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
			appRepo.Bootstrap(&testCase.routes, func() {
				if app.GetConfigValue("ENV") == gin.ReleaseMode {
					return
				}
				testFuncInsideBootstrap = "test"
			})

		}

		routesFound := 0
		if app.GetConfigValue("ENV") == gin.ReleaseMode {
			suite.Empty(testFuncInsideBootstrap)
		} else {
			suite.Equal("test", testFuncInsideBootstrap)
		}
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
