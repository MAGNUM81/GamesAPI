package api

import (
	"GamesAPI/src/database"
	"GamesAPI/src/domain"
	"GamesAPI/src/middleware"
	"GamesAPI/src/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Bootstrap(r *gin.Engine) {
	var dbInstance, dbErr = database.Setup(domain.InitRepositories)
	if dbErr != nil {
		//means that the database connection could not be created.
		//TODO: idea: add middleware here to trigger DB connection before any call
		//			  that middleware would try to connect to the DB and return 500 if it cannot,
		//			  effectively blocking any other request if DB is down.
		panic(fmt.Errorf("database connection could not be instantiated %s", dbErr.Error()))
	}
	defer dbInstance.Close()

	middleware.InitApiToken(r)
	router.InitRoutes(r)

	err := r.Run()
	HandleErrors(err)
}

func HandleErrors(err error) {
	if err != nil {
		panic("Something went horribly wrong! " + err.Error())
	}
}
