package router

import (
	"GamesAPI/src/controllers"
	"github.com/gin-gonic/gin"
)

func InitAllGameRoutes(r *gin.Engine) {
	g := InitGameRouterGroup(r)
	InitGetAllGamesRoute(g)
	InitGetGameRoute(g)
	InitCreateGameRoute(g)
	InitUpdateGameRoute(g)
	InitDeleteGameRoute(g)
}

func InitGameRouterGroup(r *gin.Engine) *gin.RouterGroup {
	return r.Group("/games")
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
