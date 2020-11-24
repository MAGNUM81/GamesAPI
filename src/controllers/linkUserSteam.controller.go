package controllers

import (
	"GamesAPI/src/External/Steam"
	"GamesAPI/src/services"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)

type input struct{
	Userid uint64 `json:"userid"`
	ProfileUrl string `json:"profile_url"`
}

func LinkSteamUser(c *gin.Context){
	i := input{}
	//json.Unmarshal(c.Request.Body, &i)
	err := c.ShouldBindJSON(&i)

	if err != nil {
		//ErrorMessageTypeCode()
		return
	}

	steamUrl := i.ProfileUrl

	if !validateUrl(steamUrl){

		//ErrorMessageTypeCode()
		return
	}

	if strings.Contains(steamUrl,"/profiles/"){
		userSteamId := strings.Split(steamUrl, "/")[4]
	}
	if strings.Contains(steamUrl,"/id/"){
		userSteamId, error := Steam.ExternalSteamUserService.GetUserID(strings.Split(steamUrl, "/")[4])

		if error != nil{
			//ErrorMessageTypeCode(c, 400, "")
			return
		}
	}
	//
	userid := i.Userid
	user,err := services.UsersService.GetUser(userid)
	if err != nil {
		//ErrorMessageTypeCode
		return
	}
	user.SteamUserID = userSteamId
	_, errorUpdate := services.UsersService.UpdateUser(user)
	if errorUpdate != nil{
		//ErrorMessageTypeCode
		return
	}
	c.JSON(200,gin.H{"Message":"Succes"})
}

func ErrorMessageTypeCode( c *gin.Context, code int, message  string) {
	c.AbortWithStatusJSON(code, gin.H{"Error": message})
}
func validateUrl(url  string) bool {
//TODO
	return true
}