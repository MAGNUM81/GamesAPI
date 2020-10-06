package router

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	//Make sure to init all routes for all modules here
	InitHomeRoutes(r)
	InitAllGameRoutes(r)
}
