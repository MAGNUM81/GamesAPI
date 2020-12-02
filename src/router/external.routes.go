package router

import (
	"GamesAPI/src/controllers"
	"github.com/gin-gonic/gin"
)

func InitExternalRoutes(group *gin.RouterGroup) {
	//Init all routes that make external calls here
	group.POST("/SyncGames", controllers.SyncGamesHandler)
	group.POST("/LinkSteamUser", controllers.LinkSteamUser)
}