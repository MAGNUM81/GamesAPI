package games

import (
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
