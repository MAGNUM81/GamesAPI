package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
	"time"
)

type Game struct {
	gorm.Model
	Title       string    `json:"title"`
	Developer   string    `json:"developer"`
	Publisher   string    `json:"publisher"`
	ReleaseDate time.Time `json:"releaseDate"`
}

func (g *Game) Validate() errorUtils.GameError {
	//TODO: add validation rules for a Game
	//		for example, we could validate that the title is not empty or nil
	return nil
}