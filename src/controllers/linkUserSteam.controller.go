package controllers

import (
	"GamesAPI/src/External/Steam"
	"GamesAPI/src/services"
	"github.com/gin-gonic/gin"
	"strings"
)

type input struct {
	Userid     uint64 `json:"userid"`
	ProfileUrl string `json:"profile_url"`
}

func LinkSteamUser(c *gin.Context) {
	i := input{}
	err := c.ShouldBindJSON(&i)

	if err != nil {
		ErrorMessageTypeCode(c, 400, "Invalid Json body")
		return
	}

	steamUrl := i.ProfileUrl

	if !validateUrl(steamUrl) {
		ErrorMessageTypeCode(c, 400, "Steam Url invalid")

		return
	}
	var userSteamId string
	var err2 error
	if strings.Contains(steamUrl, "/profiles/") {
		userSteamId = strings.Split(steamUrl, "/")[4]
	}
	if strings.Contains(steamUrl, "/id/") {
		userSteamId, err2 = Steam.ExternalSteamUserService.GetUserID(strings.Split(steamUrl, "/")[4])

		if err2 != nil {
			ErrorMessageTypeCode(c, 500, "Could not get the user Steam id from Steam Url")
			return
		}
	}
	userid := i.Userid
	user, errget := services.UsersService.GetUser(userid)
	if errget != nil {
		ErrorMessageTypeCode(c, errget.Status(), errget.Message())
		return
	}
	user.SteamUserId = userSteamId
	_, errorUpdate := services.UsersService.UpdateUser(user)
	if errorUpdate != nil {
		ErrorMessageTypeCode(c, errorUpdate.Status(), errorUpdate.Message())
		return
	}
	c.JSON(200, gin.H{"Message": "Success"})
}

func ErrorMessageTypeCode(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"Error": message})
}

func validateUrl(url string) bool {
	if strings.Contains(url, "steamcommunity") {
		if len(strings.Split(url, "/")) >= 5 {
			if strings.Split(url, "/")[3] == "profiles" || strings.Split(url, "/")[3] == "id" {
				if strings.Split(url, "/")[4] != "" {
					return true
				}
			}
		}
	}
	return false
}
