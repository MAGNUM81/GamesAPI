package router

import (
	"GamesAPI/src/controllers"
	"github.com/gin-gonic/gin"
)

func InitLoginRoute(g *gin.RouterGroup) {
	g.GET("/login", controllers.LoginController)
}

/*
func InitRefreshRoute(g *gin.RouterGroup) {
	g.GET("/refresh", controllers.RefreshController)
}
*/
