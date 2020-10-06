package controllers

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getGameId(gameIdParam string) (uint64, errorUtils.GameError) {
	gameId, gameError := strconv.ParseUint(gameIdParam, 10, 64)
	if gameError != nil {
		return 0, errorUtils.NewBadRequestError("game id should be a number")
	}
	return gameId, nil
}

func isGameError(c *gin.Context, e errorUtils.GameError) bool {
	if e != nil {
		c.JSON(e.Status(), e)
	}
	return e != nil
}

func GetGame(c *gin.Context) {
	gameId, gameErr := getGameId(c.Param("id"))
	if isGameError(c, gameErr) {
		return
	}

	game, err := services.GamesService.GetGame(gameId)
	if isGameError(c, err){
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": game})
}

func GetAllGames(c *gin.Context) {
	games, err := services.GamesService.GetAllGames()
	if isGameError(c, err){
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": games})
}

func CreateGame(c *gin.Context) {
	var game domain.Game
	if err := c.ShouldBindJSON(&game); err != nil {
		gameErr := errorUtils.NewUnprocessableEntityError("invalid json body")
		c.JSON(gameErr.Status(), gameErr)
		return
	}

	g, err := services.GamesService.CreateGame(&game)
	if isGameError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": g})
}

func UpdateGame(c *gin.Context) {
	gameId, err := getGameId(c.Param("id"))
	if isGameError(c, err) {
		return
	}

	var game domain.Game
	if err := c.ShouldBindJSON(&game); err != nil {
		gameErr := errorUtils.NewUnprocessableEntityError("invalid json body")
		c.JSON(gameErr.Status(), gameErr)
		return
	}
	game.ID = uint(gameId)
	g, err := services.GamesService.UpdateGame(&game)
	if isGameError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": g})
}

func DeleteGame(c *gin.Context) {
	gameId, err := getGameId(c.Param("id"))
	if isGameError(c, err) {
		return
	}
	if err := services.GamesService.DeleteGame(gameId); isGameError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
