package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"time"
)

type Game struct {
	ID        uint64 `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	Title       string    `json:"title"`
	Developer   string    `json:"developer"`
	Publisher   string    `json:"publisher"`
	ReleaseDate time.Time `gorm:"column:releaseDate" json:"releaseDate"`
}

func (g *Game) Validate() errorUtils.EntityError {
	//TODO: add validation rules for a Game
	//		for example, we could validate that the title is not empty or nil
	return nil
}