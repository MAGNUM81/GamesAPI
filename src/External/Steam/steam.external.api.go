package Steam

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	ExternalSteamUserService ExternalSteamUserServiceInterface = &externalSteamUserService{}
)

type externalSteamUserService struct {}


type ExternalSteamUserServiceInterface interface {
	GetUserID(personalURL string) (string, error)
	GetUserOwnedGames(userID string) ([]string, error)
}

func getFromSteam(requestURL string) ([]byte, error){
	resp, err := http.Get(requestURL)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bodyBytes, err
}

func (e externalSteamUserService) GetUserID(personalURL string) (string, error) {
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

func (e externalSteamUserService) GetUserOwnedGames(userID string) ([]string, error){
	type game struct {
		Appid                    int `json:"appid"`
		Playtime_forever         int `json:"playtime_forever"`
		Playtime_windows_forever int `json:"playtime_windows_forever"`
		Playtime_mac_forever     int `json:"playtime_mac_forever"`
		Playtime_linux_forever   int `json:"playtime_linux_forever"`
	}
	type ownedGames struct {
		Response struct {
			Game_count int `json:"game_count"`
			Games []game `json:"games"`
		}
	}
	key := os.Getenv("STEAMKEY")
	ownedGamesInfo, err := getFromSteam("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key="+ key +"&steamid="+ userID +"&format=json")
	if err != nil {
		return []string{""}, err
	}

	var userOwnedGames ownedGames
	json.Unmarshal(ownedGamesInfo, &userOwnedGames)

	var usableSteamGameIDs []string;
	for _, games := range userOwnedGames.Response.Games {
		usableSteamGameIDs = append(usableSteamGameIDs, strconv.Itoa(games.Appid))
	}

	return usableSteamGameIDs, nil
}
