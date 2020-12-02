package Steam

import (
	"GamesAPI/src/domain"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ExternalSteamUserService ExternalSteamUserServiceInterface = &externalSteamUserService{}
)

type externalSteamUserService struct{}

type ExternalSteamUserServiceInterface interface {
	GetUserID(personalURL string) (string, error)
	GetUserOwnedGames(userID string) ([]string, error)
	GetGameInfo(gameID string) (domain.Game, error)
}

func getFromSteam(requestURL string) ([]byte, error) {
	resp, err := http.Get(requestURL)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = resp.Body.Close()

	return bodyBytes, err
}

func (e externalSteamUserService) GetUserID(personalURL string) (string, error) {
	key := os.Getenv("STEAMKEY")
	steamID, err := getFromSteam("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=" + key + "&vanityurl=" + personalURL)
	if err != nil {
		return "", err
	}
	var userinfo basicUserSteamType
	json.Unmarshal(steamID, &userinfo)

	if userinfo.Response.Success == 1 {
		return userinfo.Response.Steamid, nil
	} else {
		return "", errors.New("no match found")
	}
}

func (e externalSteamUserService) GetUserOwnedGames(userID string) ([]string, error){
	key := os.Getenv("STEAMKEY")
	ownedGamesInfo, err := getFromSteam("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=" + key + "&steamid=" + userID + "&format=json")
	if err != nil {
		return []string{""}, err
	}

	var userOwnedGames ownedGamesSteamType
	json.Unmarshal(ownedGamesInfo, &userOwnedGames)

	var usableSteamGameIDs []string
	for _, games := range userOwnedGames.Response.Games {
		usableSteamGameIDs = append(usableSteamGameIDs, strconv.Itoa(games.Appid))
	}

	return usableSteamGameIDs, nil
}

func (e externalSteamUserService) GetGameInfo(gameID string) (domain.Game, error){
	gameInfo, err := getFromSteam("https://store.steampowered.com/api/appdetails?appids="+gameID)
	if err != nil {
		return domain.Game{}, err
	}

	var unFilterGameInfo map[string]interface{}
	err = json.Unmarshal(gameInfo, &unFilterGameInfo)
	if err != nil {
		return domain.Game{}, err
	}
	var successMaybeSomeday = unFilterGameInfo[gameID]
	if successMaybeSomeday == nil {
		return domain.Game{}, errors.New("server did not respond correctly")
	}
	if !unFilterGameInfo[gameID].(map[string]interface{})["success"].(bool) {
		return domain.Game{}, errors.New("bad game ID")
	}

	data := unFilterGameInfo[gameID].(map[string]interface{})["data"]
	if data == nil {
		return domain.Game{}, errors.New("game data is invalid")
	}
	//Devilish line used to get 2 level deeper into the json struct without knowing the exact way the struct is defined
	var pureGameInfo = data.(map[string]interface{})
	var filteredGameInfo domain.Game

	devs := pureGameInfo["developers"]
	if devs == nil {
		filteredGameInfo.Developer = ""
	} else {
		filteredGameInfo.Developer = interfaceListToSingleString(pureGameInfo["developers"].([]interface{}))
	}
	pubs := pureGameInfo["publishers"]
	if pubs == nil {
		filteredGameInfo.Publisher = ""
	} else {
		filteredGameInfo.Publisher = interfaceListToSingleString(pureGameInfo["publishers"].([]interface{}))
	}
	filteredGameInfo.Title = pureGameInfo["name"].(string)
	filteredGameInfo.SteamId = gameID

	//var date = pureGameInfo["release_date"].(map[string]interface{})
	//year,month,day := steamTimeParser(date["date"].(string))
	//filteredGameInfo.ReleaseDate = time.Date(year,month,day,0,0,0,0,time.UTC)
	return filteredGameInfo, nil
}

func interfaceListToSingleString(interfaceList []interface{}) string {
	var singleString string
	for x := range interfaceList {
		singleString += interfaceList[x].(string) + " | "
	}
	//will remove the last useless separator
	return strings.TrimSuffix(singleString, " | ")
}

func steamTimeParser(steamTime string)(int, time.Month, int) {
	splitedTime := strings.Split(steamTime, " ")
	day,_ := strconv.Atoi(splitedTime[0])
	year,_ := strconv.Atoi(splitedTime[2])
	switch splitedTime[1] {
	case "Jan,":
		return year, 1, day
	case "Feb,":
		return year, 2, day
	case "Mar,":
		return year, 3, day
	case "Apr,":
		return year, 4, day
	case "May,":
		return year, 5, day
	case "Jun,":
		return year, 6, day
	case "Jul,":
		return year, 7, day
	case "Aug,":
		return year, 8, day
	case "Sep,":
		return year, 9, day
	case "Oct,":
		return year, 10, day
	case "Nov,":
		return year, 11, day
	default :
		return year, 12, day
	}
}
