package router

import (
	"GamesAPI/src/controllers"
	"github.com/gin-gonic/gin"
)

func InitHomeRoutes(g *gin.RouterGroup) {
	g.GET("/", controllers.Home)
}
