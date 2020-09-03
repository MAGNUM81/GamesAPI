package games

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Game struct {
	gorm.Model
	Id          uint      `json:"id" gorm:"primary_key"`
	Title       string    `json:"title"`
	Developer   string    `json:"developer"`
	Publisher   string    `json:"publisher"`
	ReleaseDate time.Time `json:"releaseDate"`
}

type CreateGameInput struct {
	Title       string `json:"title" binding:"required"`
	Developer   string `json:"developer" binding:"required"`
	Publisher   string `json:"publisher" binding:"required"`
	ReleaseDate string `json:"releaseDate" binding:"required"`
}

type UpdateGameInput struct {
	Title       string `json:"title"`
	Developer   string `json:"developer"`
	Publisher   string `json:"publisher"`
	ReleaseDate string `json:"releaseDate"`
}
