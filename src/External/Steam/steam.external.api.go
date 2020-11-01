package Steam

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)


func getFromSteam(requestURL string) ([]byte, error){
	resp, err := http.Get(requestURL)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bodyBytes, err
}

func GetUserID(personalURL string) (string, error) {
	type basicUser struct {
		Response struct {
			Steamid string `json:"steamid"`
			Success int    `json:"success"`
		} `json:"response"`
	}
	key := os.Getenv("STEAMKEY")
	steamID, err := getFromSteam("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=" + key + "&vanityurl=" + personalURL)
	if err != nil {
		return "", err
	}
	var userinfo basicUser
	json.Unmarshal(steamID, &userinfo)

	if userinfo.Response.Success == 1{
		return userinfo.Response.Steamid, nil
	} else {
		return "",  errors.New("no match found")
	}
}

