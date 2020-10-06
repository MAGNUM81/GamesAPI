package router

import (
	"GamesAPI/src/controllers"
	"github.com/gin-gonic/gin"
)

func InitHomeRoutes(r *gin.Engine) {
	r.GET("/", controllers.Home)
}

