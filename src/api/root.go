package api

import (
	"GamesAPI/src/database"
	"GamesAPI/src/router"
	"github.com/gin-gonic/gin"
)

func Bootstrap(r *gin.Engine) {
	ConnectDatabase()
	defer database.Instance.Close()
	router.InitRoutes(r)

	err := r.Run()
	HandleErrors(err)

}

func HandleErrors(err error) {
	if err != nil {
		panic("Something went horribly wrong!")
	}
}
