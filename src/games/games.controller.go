package games

import (
	"GamesAPI/src/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func FindGames(c *gin.Context) {
	var games []Game
	database.Instance.Find(&games)

	c.JSON(http.StatusOK, gin.H{"data": games})
}

func FindGame(c *gin.Context) {
	var game Game
	if err := database.Instance.Where("id = ?", c.Param("id")).First(&game).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": game})
}

func CreateGame(c *gin.Context) {
	var input CreateGameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	releaseDate, err := time.Parse("2006-01-02", input.ReleaseDate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The date wasn't formatted properly!"})
		return
	}

	game := Game{
		Title:       input.Title,
		Developer:   input.Developer,
		Publisher:   input.Publisher,
		ReleaseDate: releaseDate}
	database.Instance.Create(&game)

	c.JSON(http.StatusOK, gin.H{"data": game})
}

func UpdateGame(c *gin.Context) {
	var game Game
	if err := database.Instance.Where("id = ?", c.Param("id")).First(&game).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var input UpdateGameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Instance.Model(&game).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": game})
}

func DeleteGame(c *gin.Context) {
	var game Game

	if err := database.Instance.Where("id = ?", c.Param("id")).First(&game).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	database.Instance.Delete(&game) //FIXME: does not set the deleted_at property in the Instance

	c.JSON(http.StatusOK, gin.H{"data": true})
}
