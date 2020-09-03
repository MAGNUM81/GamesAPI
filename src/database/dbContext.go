package database

import (
	"GamesAPI/src/games"
)

type IDbContext interface {
	New() DbContext
	Find(id interface{}) (interface{}, error)
	SaveChanges() error
}

type DbContext struct {

}

func (DbContext) New(f func(...IDbContext) IDbContext) IDbContext {
	
}

func (ctx *GamesDbContext) Find(id uint) (*games.Game, error) {

}

func (ctx *GamesDbContext) SaveChanges() error {

}