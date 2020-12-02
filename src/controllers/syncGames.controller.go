package controllers

import (
	"GamesAPI/src/External/Steam"
	"GamesAPI/src/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type inputSyncGames struct{
	Userid uint64 `json:"userid"`
}
/*	1. 	trouver userid dans json, parse en uint64 -> si err , bad request
	2. 	aller chercher le user associé en BD -> si pas trouvé, bad request
	3. 	extraire steamId -> si vide, not found + message custom
	4. 	aller chercher tous les game ids -> si err, 500 + err
	5. 	pour chacun des game ids :
		.1 si un jeu n'existe pas avec game id :
			obtenir le jeu steam associé (domain.Game) -> si erreur, 500
			créer le jeu en BD -> si erreur, forward error
	6. 	200
*/
func SyncGamesHandler(c *gin.Context) {
	input := inputSyncGames{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		AbortWithStatusError(c, http.StatusBadRequest, err)
		return
	}

	user, errGetUser := services.UsersService.GetUser(input.Userid)
	if errGetUser != nil {
		AbortWithStatusError(c, errGetUser.Status(), errGetUser)
		return
	}

	steamUserId := user.SteamUserId
	if steamUserId == "" {
		AbortWithStatusError(c, http.StatusNotFound, errors.New("l'usager ne possède pas de ID steam"))
		return
	}

	gameIds, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserId)
	if err != nil {
		AbortWithStatusError(c, http.StatusInternalServerError, err)
	}

	for _, gameId := range gameIds {
		existsGameWithSteamId, errExists := services.GamesService.ExistsWithSteamID(gameId)
		if errExists != nil {
			AbortWithStatusError(c, errExists.Status(), errExists)
			return
		}
		if !existsGameWithSteamId {
			g, err := Steam.ExternalSteamUserService.GetGameInfo(gameId)
			if err != nil {
				AbortWithStatusError(c, http.StatusInternalServerError, err)
				return
			}
			_, errCreate := services.GamesService.CreateGame(&g)
			if errCreate != nil {
				AbortWithStatusError(c, errCreate.Status(), errCreate)
				return
			}
		}
		c.JSON(http.StatusOK, nil)
	}
}

func AbortWithStatusError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, err)
}