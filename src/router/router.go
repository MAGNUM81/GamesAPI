package router

import (
	"GamesAPI/src/middleware"
	"github.com/gin-gonic/gin"
)

func InitAllRoutes(r *gin.Engine) {

	coreGroup := initCoreGroup(r)
	authGroup := initAuthGroup(r)

	middleware.InitApiToken(coreGroup, authGroup) //will apply to all routes

	middleware.InitAuthorization(coreGroup)
	middleware.InitUserSessionHandler(coreGroup)
}

func initCoreGroup(r *gin.Engine) *gin.RouterGroup {
	//Make sure to init all routes for all modules here
	g := r.Group("")
	InitHomeRoutes(g)
	InitAllGameRoutes(g)
	InitAllUserRoutes(g)
	return g
}

func initAuthGroup(r *gin.Engine) *gin.RouterGroup {
	g := r.Group("/auth")
	InitLoginRoute(g)
	//InitRefreshRoute(g)
	return g
}
