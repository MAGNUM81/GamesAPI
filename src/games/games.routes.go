package games

import "github.com/gin-gonic/gin"

func InitRoutes(r *gin.Engine) {
	g := r.Group("/games")
	g.GET("", FindGames)
	g.GET("/:id", FindGame) // new
	g.POST("", CreateGame)
	g.PATCH("/:id", UpdateGame)
	g.DELETE("/:id", DeleteGame)
}
