package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"Message": "Welcome to GamesAPI!"})
}
