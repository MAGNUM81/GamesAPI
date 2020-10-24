package controllers

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getGameId(gameIdParam string) (uint64, errorUtils.EntityError) {
	gameId, gameError := strconv.ParseUint(gameIdParam, 10, 64)
	if gameError != nil {
		return 0, errorUtils.NewBadRequestError("game id should be a number")
	}
	return gameId, nil
}

func GetGame(c *gin.Context) {
	gameId, gameErr := getGameId(c.Param("id"))
	if errorUtils.IsEntityError(c, gameErr) {
		return
	}

	game, err := services.GamesService.GetGame(gameId)
	if errorUtils.IsEntityError(c, err){
		return
	}

	c.JSON(http.StatusOK, game)
}

func GetAllGames(c *gin.Context) {
	games, err := services.GamesService.GetAllGames()
	if errorUtils.IsEntityError(c, err){
		return
	}

	c.JSON(http.StatusOK, games)
}

func CreateGame(c *gin.Context) {
	var game domain.Game
	if err := c.ShouldBindJSON(&game); game.Validate() != nil || err != nil {
		gameErr := errorUtils.NewUnprocessableEntityError("invalid json body")
		c.JSON(gameErr.Status(), gameErr)
		return
	}

	g, err := services.GamesService.CreateGame(&game)
	if errorUtils.IsEntityError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, g)
}

func UpdateGame(c *gin.Context) {
	gameId, err := getGameId(c.Param("id"))
	if errorUtils.IsEntityError(c, err) {
		return
	}

	var game domain.Game
	if err := c.ShouldBindJSON(&game); game.Validate() != nil || err != nil {
		gameErr := errorUtils.NewUnprocessableEntityError("invalid json body")
		c.JSON(gameErr.Status(), gameErr)
		return
	}
	game.ID = gameId
	g, err := services.GamesService.UpdateGame(&game)
	if errorUtils.IsEntityError(c, err) {
		return
	}

	c.JSON(http.StatusOK, g)
}

func DeleteGame(c *gin.Context) {
	gameId, err := getGameId(c.Param("id"))
	if errorUtils.IsEntityError(c, err) {
		return
	}
	if err := services.GamesService.DeleteGame(gameId); errorUtils.IsEntityError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
