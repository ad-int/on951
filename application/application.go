package application

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"github.com/pascaldekloe/jwt"
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

var application *TApplication

type TApplication struct {
	config       map[string]string
	router       router.TAppRouter
	articlesRepo *api.TArticlesRepository
}

func GetApplication() *TApplication {
	if application == nil {
		application = &TApplication{}
		return application
	}
	return application
}

func Bootstrap(routesList *map[string]router.TRoutesList, userFuncs ...func()) {
	appInstance := GetApplication()
	appInstance.ReadEnvFile()
	appInstance.InitDb()
	for _, userFunc := range userFuncs {
		userFunc()
	}
	appInstance.Init(routesList)
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
	db := database.TDatabase{}
	connOk := db.ConnectToDB(GetApplication().GetConfigValue("DSN"))
	if !connOk {
		log.Fatalln("error connecting to db")
	}
	app.articlesRepo = &api.TArticlesRepository{
		TDatabase: db,
	}
	app.GetArticlesRepo().AutoMigrate()
	return connOk
}
func (app *TApplication) Init(routes *map[string]router.TRoutesList) error {

	if len(GetImagesDir()) < 1 {
		log.Fatalln("Cannot read images directory path")
	}
	log.Println(GetImagesDir())

	app.router.Configure()
	app.router.InitRoutes(routes)
	return app.router.Run()

}

func GetAuthorizedUserFromHeader(authHeader string) (dbStructure.User, error) {

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

	claim, err := jwt.HMACCheck([]byte(authHeaderParts[1]), []byte(GetApplication().GetConfigValue("SECRET")))
	if err != nil {
		return dbStructure.User{}, err
	}
	if !claim.Valid(time.Now()) {
		return dbStructure.User{}, errors.New(MsgUnauthorized)
	}

	if !claim.AcceptAudience(GetApplication().GetConfigValue("AUDIENCE")) {
		return dbStructure.User{}, errors.New(MsgNotAcceptedAudience)
	}

	if claim.Issuer != GetApplication().GetConfigValue("ISSUER") {
		return dbStructure.User{}, errors.New(MsgInvalidIssuer)
	}
	var user dbStructure.User
	err = json.Unmarshal([]byte(claim.Subject), &user)

	return user, err
}

func GetImagesDir() string {
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
