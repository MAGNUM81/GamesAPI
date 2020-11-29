package router

import (
	"GamesAPI/src/controllers"
	"github.com/gin-gonic/gin"
)

func InitAllGameRoutes(root *gin.RouterGroup) {
	g := InitGameRouterGroup(root)
	InitGetAllGamesRoute(g)
	InitGetGameRoute(g)
	InitCreateGameRoute(g)
	InitUpdateGameRoute(g)
	InitDeleteGameRoute(g)
}

func InitGameRouterGroup(g *gin.RouterGroup) *gin.RouterGroup {
	return g.Group("/games")
}

func InitGetAllGamesRoute(g *gin.RouterGroup) {
	g.GET("", controllers.GetAllGames)
}

func InitGetGameRoute(g *gin.RouterGroup) {
	g.GET("/:id", controllers.GetGame)
}

func InitCreateGameRoute(g *gin.RouterGroup) {
	g.POST("", controllers.CreateGame)
}

func InitUpdateGameRoute(g *gin.RouterGroup) {
	g.PATCH("/:id", controllers.UpdateGame)
}

func InitDeleteGameRoute(g *gin.RouterGroup) {
	g.DELETE("/:id", controllers.DeleteGame)
}
