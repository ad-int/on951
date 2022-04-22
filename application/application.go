package application

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"github.com/pascaldekloe/jwt"
	"io/fs"
	"log"
	"on951/api"
	"on951/database"
	dbStructure "on951/database/structure"
	"on951/router"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MsgUnauthorized = "unauthorized"
const MsgNotAcceptedAudience = "not accepted audience"
const MsgNoAuthorizationToken = "missing authorization token"
const MsgInvalidAuthorizationToken = "invalid authorization token"
const MsgInvalidIssuer = "invalid issuer"

const ImagesDirectory = "images"

type IApplicationRepository interface {
	GetApplication() *TApplication
	Bootstrap(routesList *map[string]router.TRoutesList, userFuncs ...func())
}

type IApplication interface {
	GetConfigValue(key string) string
	GetImagesDir() string
	InitImagesDir() string
	InitDb() bool
	ReadEnvFile() bool
	Init(routes *map[string]router.TRoutesList) error
	GetAuthorizedUserFromHeader(authHeader string) (dbStructure.User, error)
	GetArticlesRepo() *api.TArticlesRepository
	GetRouter() router.AppRouter
}

var applicationRepository = &TApplicationRepository{}

type TApplicationRepository struct {
	application IApplication
}

type TApplication struct {
	config       map[string]string
	router       router.TAppRouter
	articlesRepo *api.TArticlesRepository
}

func (applicationRepository *TApplicationRepository) GetApplication() IApplication {
	if applicationRepository.application == nil {
		applicationRepository.application = &TApplication{}
	}
	return applicationRepository.application
}

func GetApplicationRepository() *TApplicationRepository {
	if applicationRepository == nil {
		applicationRepository = &TApplicationRepository{}
	}
	return applicationRepository
}
func GetApplication() IApplication {
	return applicationRepository.GetApplication()
}

func (applicationRepository *TApplicationRepository) Bootstrap(routesList *map[string]router.TRoutesList, userFuncs ...func()) {
	appInstance := applicationRepository.GetApplication()
	appInstance.ReadEnvFile()

	appInstance.InitDb()
	for _, userFunc := range userFuncs {
		userFunc()
	}
	if err := appInstance.Init(routesList); err != nil {
		panic(err)
	}
}

func (app *TApplication) GetArticlesRepo() *api.TArticlesRepository {
	return app.articlesRepo
}

func (app *TApplication) GetConfigValue(key string) string {
	return app.config[key]
}

func (app *TApplication) ReadEnvFile() bool {
	var err error
	app.config, err = godotenv.Read(".env")
	if err != nil {
		log.Println(err)
		log.Fatalln("Unable to read .env file!")
		return false
	}
	return true
}

func (app *TApplication) InitDb() bool {
	db := &database.TDatabase{}
	connOk := db.ConnectToDB(app.GetConfigValue("DSN"))
	if !connOk {
		log.Fatalln("error connecting to db")
	}
	app.articlesRepo = &api.TArticlesRepository{
		IDatabase: db,
	}
	app.GetArticlesRepo().AutoMigrate()
	return connOk
}
func (app *TApplication) Init(routes *map[string]router.TRoutesList) error {

	if len(app.GetImagesDir()) < 1 {
		log.Fatalln("Cannot read images directory path")
	}
	log.Println(app.GetImagesDir())

	app.router.Configure()
	app.router.InitRoutes(routes)
	return app.router.Run()

}

func (app *TApplication) GetAuthorizedUserFromHeader(authHeader string) (dbStructure.User, error) {

	if authHeader == "" {
		return dbStructure.User{}, errors.New(MsgUnauthorized)
	}
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 {
		return dbStructure.User{}, errors.New(MsgNoAuthorizationToken)
	}

	if authHeaderParts[0] != "Bearer" {
		return dbStructure.User{}, errors.New(MsgInvalidAuthorizationToken)
	}

	claim, err := jwt.HMACCheck([]byte(authHeaderParts[1]), []byte(app.GetConfigValue("SECRET")))
	if err != nil {
		return dbStructure.User{}, err
	}
	if !claim.Valid(time.Now()) {
		return dbStructure.User{}, errors.New(MsgUnauthorized)
	}

	if !claim.AcceptAudience(app.GetConfigValue("AUDIENCE")) {
		return dbStructure.User{}, errors.New(MsgNotAcceptedAudience)
	}

	if claim.Issuer != app.GetConfigValue("ISSUER") {
		return dbStructure.User{}, errors.New(MsgInvalidIssuer)
	}
	var user dbStructure.User
	err = json.Unmarshal([]byte(claim.Subject), &user)

	return user, err
}

func (app *TApplication) GetImagesDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return ""
	}
	iDir := filepath.Join(dir, ImagesDirectory)
	fi, err := os.Stat(iDir)
	if err != nil {
		log.Println(err)
		return ""
	}
	if !fi.IsDir() {
		return ""
	}
	return iDir
}
func (app *TApplication) InitImagesDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return ""
	}
	iDir := filepath.Join(dir, ImagesDirectory)
	err = os.Mkdir(iDir, fs.ModeDir)
	fi, err := os.Stat(iDir)
	if err != nil {
		return ""
	}
	if !fi.IsDir() {
		return ""
	}
	return iDir
}

func (app *TApplication) GetRouter() router.AppRouter {
	return &app.router
}
