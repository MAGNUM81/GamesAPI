package controllers

import (
	"GamesAPI/src/External/Steam"
	"GamesAPI/src/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

//TODO
//Error code to change !

type input struct{
	Userid uint64 `json:"userid"`
	ProfileUrl string `json:"profile_url"`
}

func LinkSteamUser(c *gin.Context){
	i := input{}
	err := c.ShouldBindJSON(&i)

	if err != nil {
		ErrorMessageTypeCode(c, 400, "Invalid Json body")
		return
	}

	steamUrl := i.ProfileUrl

	if !validateUrl(steamUrl){
		ErrorMessageTypeCode(c, 400, "Steam Url invalid")

		return
	}
	if strings.Contains(steamUrl,"/profiles/"){
		userSteamId := strings.Split(steamUrl, "/")[4]
	}
	if strings.Contains(steamUrl,"/id/"){
		userSteamId, error := Steam.ExternalSteamUserService.GetUserID(strings.Split(steamUrl, "/")[4])

		if error != nil{
			ErrorMessageTypeCode(c, 400, "Could not get the user Steam id from Steam Url")
			return
		}
	}
	userid := i.Userid
	user,err := services.UsersService.GetUser(userid)
	if err != nil {
		ErrorMessageTypeCode(c, 400, "Service User could not get user from User id" )
		return
	}
	user.SteamUserID = userSteamId
	_, errorUpdate := services.UsersService.UpdateUser(user)
	if errorUpdate != nil{
		ErrorMessageTypeCode(c, 400,"Error in service when updated User")
		return
	}
	c.JSON(200,gin.H{"Message":"Succes"})
}

func ErrorMessageTypeCode( c *gin.Context, code int, message  string) {
	c.AbortWithStatusJSON(code, gin.H{"Error": message})
}

func validateUrl(url  string) bool {
//TODO
	//url similar to  https://steamcommunity.com/profiles/############
	//url doit etre https -- !?!
	//steamcommunity
	//split aumoins de 5
	// split position 3 est /profile/ ou /id/
	if strings.Contains(url ,"steamcommunity"){
		if len(strings.Split(url, "/")) >= 5{
			if strings.Split(url, "/")[3] == "profiles" || strings.Split(url, "/")[3] == "id" {
				return true
		}}}
	return false
}