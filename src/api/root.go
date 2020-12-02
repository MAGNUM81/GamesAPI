package api

import (
	"GamesAPI/src/database"
	"GamesAPI/src/domain"
	"GamesAPI/src/router"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/authUtils"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
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

	//let's add a Session that doesn't expire for devs
	//NOT FOR PROD
	sessionKey := "2837503506"
	_, _ = services.UserSessionService.CreateSession(&domain.UserSession{
		Token:     sessionKey,
		UserId:    1,
		ExpiresAt: time.Now().AddDate(1, 0, 0).UnixNano(), //token will expire 1 year after server boot up
	})
	h, _ := authUtils.HashAndSalt([]byte("network7"))
	users, _ := services.UsersService.GetAllUsers()
	_ = fmt.Sprintf("%v", users)
	u, _ := services.UsersService.GetUser(uint64(1))
	masterEmail := "master@test.com"
	if u == nil {
		_, _ = services.UsersService.CreateUser(&domain.User{
			Name:         "master",
			Email:        masterEmail,
			PasswordHash: h,
		})
	}
	role, _ := services.UserRoleService.GetRolesByUserID(uint64(1))
	if len(role) == 0 {
		_, _ = services.UserRoleService.CreateRole(&domain.UserRole{
			UserID: 1,
			Name:   "admin",
		})
	}
	_, _ = fmt.Printf("This is the bypass session key : %s\n", sessionKey)
	_, _ = fmt.Printf("This is the master email : %s\n", masterEmail)
	//END : NOT FOR PROD

	router.InitAllRoutes(r)

	err := r.Run()
	HandleErrors(err)
}

func HandleErrors(err error) {
	if err != nil {
		panic("Something went horribly wrong! " + err.Error())
	}
}
