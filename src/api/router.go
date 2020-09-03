package api

import (
	"GamesAPI/src/games"
	"GamesAPI/src/home"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	//Make sure to init all routes for all modules here
	home.InitRoutes(r)
	games.InitRoutes(r)
}
