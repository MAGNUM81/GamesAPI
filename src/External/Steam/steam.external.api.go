package Steam

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//TODO Something better than this...
const steamKey = "9230546D5E965861D940A995413DB4C8"

func getFromSteam(requestURL string) []byte{
	resp, err := http.Get(requestURL)
	if err != nil{
		log.Panic(err)
		return nil
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return bodyBytes
}

func GetUserID(personalURL string) string {
	type basicUser struct {
		Response struct {
			Steamid string `json:"steamid"`
			Success int    `json:"success"`
		} `json:"response"`
	}
	steamID := getFromSteam("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=" + steamKey + "&vanityurl=" + personalURL)
	if steamID == nil {
		return "Unexpected error"
	}
	var userinfo basicUser
	json.Unmarshal(steamID, &userinfo)

	if userinfo.Response.Success == 1{
		return userinfo.Response.Steamid
	} else {
		return "No match"
	}
}

