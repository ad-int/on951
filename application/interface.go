package application

import (
	"on951/api"
	"on951/database"
	dbStructure "on951/database/structure"
	"on951/router"
)

type IApplicationRepository interface {
	GetApplication() IApplication
	Bootstrap(routesList *map[string]router.TRoutesList, userFuncs ...func())
}

type IApplication interface {
	GetConfigValue(key string) string
	GetImagesDir() string
	InitDb() bool
	GetDatabase() database.IDatabase
	ReadEnvFile() (bool, map[string]string)
	Init(routes *map[string]router.TRoutesList) error
	GetAuthorizedUserFromHeader(authHeader string) (dbStructure.User, error)
	SetArticlesRepo(repository api.ArticlesRepository)
	GetArticlesRepo() api.ArticlesRepository
	GetRouter() router.AppRouter
}
