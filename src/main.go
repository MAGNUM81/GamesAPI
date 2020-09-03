package main

import (
	"GamesAPI/src/api"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Go Games API")
	r := gin.Default()
	api.Bootstrap(r)
}