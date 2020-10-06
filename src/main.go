package main

import (
	"GamesAPI/src/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Go Games API")
	r := gin.Default()
	api.Bootstrap(r)
}