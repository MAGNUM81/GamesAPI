package database

import (
	"GamesAPI/src/games"
	"github.com/jinzhu/gorm"
)

func BindEntities(database *gorm.DB) {
	//Make sure to call this method on each model you have created.
	database.AutoMigrate(&games.Game{})
}
