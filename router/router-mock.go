package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type TAppRouterMock struct {
	mock.Mock
	AppRouter
}

func (appRouter *TAppRouterMock) Configure() error {
	args := appRouter.Called()
	return args.Error(0)
}

func (appRouter *TAppRouterMock) InitRoutes(routes *map[string]TRoutesList) {
	_ = appRouter.Called(routes)
}

func (appRouter *TAppRouterMock) GetEngine() *gin.Engine {
	args := appRouter.Called()
	return args.Get(0).(*gin.Engine)
}

func (appRouter *TAppRouterMock) Run() error {
	args := appRouter.Called()
	return args.Error(0)
}
