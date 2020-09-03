package api

import (
	"GamesAPI/src/database"
	"github.com/gin-gonic/gin"
)

func Bootstrap(r *gin.Engine) {
	database.Instance = ConnectDatabase()
	defer database.Instance.Close()
	InitRoutes(r)

	err := r.Run()
	HandleErrors(err)
}

func HandleErrors(err error) {
	if err != nil {
		panic("Something went horribly wrong!\n" + err.Error())
	}
}
