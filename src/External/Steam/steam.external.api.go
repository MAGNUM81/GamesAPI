package Steam

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"GamesAPI/src/domain"
	"strings"
	"time"
)

var (
	ExternalSteamUserService ExternalSteamUserServiceInterface = &externalSteamUserService{}
)

type externalSteamUserService struct {}


type ExternalSteamUserServiceInterface interface {
	GetUserID(personalURL string) (string, error)
	GetUserOwnedGames(userID string) ([]string, error)
	GetGameInfo(gameID string) (domain.Game, error)
}

func getFromSteam(requestURL string) ([]byte, error){
	resp, err := http.Get(requestURL)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

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

	if userinfo.Response.Success == 1{
		return userinfo.Response.Steamid, nil
	} else {
		return "",  errors.New("no match found")
	}
}

func (e externalSteamUserService) GetUserOwnedGames(userID string) ([]string, error){
	key := os.Getenv("STEAMKEY")
	ownedGamesInfo, err := getFromSteam("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key="+ key +"&steamid="+ userID +"&format=json")
	if err != nil {
		return []string{""}, err
	}

	var userOwnedGames ownedGamesSteamType
	json.Unmarshal(ownedGamesInfo, &userOwnedGames)

	var usableSteamGameIDs []string;
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
	var unFilterGameInfo getGameInfoSteamType
	json.Unmarshal(gameInfo, &unFilterGameInfo)

	var filteredGameInfo domain.Game
	filteredGameInfo.Title = unFilterGameInfo.GameInfo.Data.Name
	filteredGameInfo.Developer = unFilterGameInfo.GameInfo.Data.Developers[0]
	filteredGameInfo.Publisher = unFilterGameInfo.GameInfo.Data.Publishers[0]
	year,month,day := steamTimeParser(unFilterGameInfo.GameInfo.Data.ReleaseDate.Date)
	filteredGameInfo.ReleaseDate = time.Date(year,month,day,0,0,0,0,nil)
	return filteredGameInfo, nil
}

func steamTimeParser(steamTime string)(int, time.Month, int) {
	splitedTime := strings.Split(steamTime, " ")
	year,_ := strconv.Atoi(splitedTime[0])
	day,_ := strconv.Atoi(splitedTime[2])
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
