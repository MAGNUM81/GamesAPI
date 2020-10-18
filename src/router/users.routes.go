package router

import (
	"GamesAPI/src/controllers"
	"github.com/gin-gonic/gin"
)

func InitAllUserRoutes(r *gin.Engine) {
	g := InitUserRouterGroup(r)
	InitGetAllUsersRoute(g, controllers.GetAllUsers)
	InitGetUserRoute(g, controllers.GetUser)
	InitCreateUserRoute(g, controllers.CreateUser)
	InitUpdateUserRoute(g, controllers.UpdateUser)
	InitDeleteUserRoute(g, controllers.DeleteUser)
}

func InitUserRouterGroup(r *gin.Engine) *gin.RouterGroup {
	return r.Group("/users")
}

func InitGetAllUsersRoute(g *gin.RouterGroup, handlerFunc gin.HandlerFunc) {
	g.GET("", handlerFunc)
}

func InitGetUserRoute(g *gin.RouterGroup, handlerFunc gin.HandlerFunc) {
	g.GET("/:id", handlerFunc)
}

func InitCreateUserRoute(g *gin.RouterGroup, handlerFunc gin.HandlerFunc) {
	g.POST("", handlerFunc)
}

func InitUpdateUserRoute(g *gin.RouterGroup, handlerFunc gin.HandlerFunc) {
	g.PATCH("/:id", handlerFunc)
}

func InitDeleteUserRoute(g *gin.RouterGroup, handlerFunc gin.HandlerFunc) {
	g.DELETE("/:id", handlerFunc)
}
