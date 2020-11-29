package router

import (
	"GamesAPI/src/middleware"
	"github.com/gin-gonic/gin"
)

func InitAllRoutes(r *gin.Engine) {

	middleware.InitApiToken(r) //will apply to all routes
	rootGroup := r.Group("")
	{
		initAuthGroup(rootGroup)
		initCoreGroup(rootGroup)
	}
}

func initCoreGroup(r *gin.RouterGroup) {
	//Make sure to init all routes for all modules here
	coreGroup := r.Group("")
	{
		middleware.InitUserSessionHandler(coreGroup)
		middleware.InitAuthorization(coreGroup)
		InitHomeRoutes(coreGroup)
		InitAllGameRoutes(coreGroup)
		InitAllUserRoutes(coreGroup)
	}
}

func initAuthGroup(g *gin.RouterGroup){
	auth := g.Group("/auth")
	{
		InitLoginRoute(auth)
		//InitRefreshRoute(g)
	}
}
