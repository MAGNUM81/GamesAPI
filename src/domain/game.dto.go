package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"time"
)

type Game struct {
	ID          uint64     `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
	Title       string     `json:"title"`
	Developer   string     `json:"developer"`
	Publisher   string     `json:"publisher"`
	ReleaseDate time.Time  `gorm:"column:releaseDate" json:"releaseDate"`
}

func (g *Game) Validate() errorUtils.EntityError {
	//check for empty title
	if g.Title == "" {
		return errorUtils.NewUnprocessableEntityError("Game title cannot be empty")
	}

	//check for empty developer
	if g.Developer == "" {
		return errorUtils.NewUnprocessableEntityError("Game developer cannot be empty")
	}

	//check for empty publisher
	if g.Publisher == "" {
		return errorUtils.NewUnprocessableEntityError("Game publisher cannot be empty")
	}
	return nil
}
