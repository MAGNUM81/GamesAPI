package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
)

var (
	GameRepo GameRepoInterface = &gameRepo{}
)

type GameRepoInterface interface {
	Get(uint64) (*Game, errorUtils.EntityError)
	Create(*Game) (*Game, errorUtils.EntityError)
	Update(*Game) (*Game, errorUtils.EntityError)
	Delete(uint64) errorUtils.EntityError
	GetAll() ([]Game, errorUtils.EntityError)
	Initialize(*gorm.DB)
}

type gameRepo struct {
	db *gorm.DB
}

func (g *gameRepo) Initialize(db *gorm.DB) {
	g.db = db
	db.AutoMigrate(&Game{})
}

func NewGameRepository(db *gorm.DB) GameRepoInterface {
	return &gameRepo{db: db}
}

func (g *gameRepo) Get(gameId uint64) (*Game, errorUtils.EntityError) {
	var game Game
	if err := g.db.Where("id = ?", gameId).First(&game).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	return &game, nil
}

func (g *gameRepo) Create(game *Game) (*Game, errorUtils.EntityError) {
	if dbc := g.db.Create(game); dbc.Error != nil {
		return nil, errorUtils.NewInternalServerError(dbc.Error.Error())
	}
	return game, nil
}

func (g *gameRepo) Update(game *Game) (*Game, errorUtils.EntityError) {
	if err := g.db.Where("id = ?", game.ID).First(&game).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	g.db.Save(*game)
	return game, nil
}

func (g *gameRepo) Delete(gameId uint64) errorUtils.EntityError {
	var game Game
	if err := g.db.Where("id = ?", gameId).First(&game).Error; err != nil {
		return errorUtils.NewNotFoundError(err.Error())
	}
	dbc := g.db.Delete(&game)
	return errorUtils.NewEntityError(dbc.Error)
}

func (g *gameRepo) GetAll() ([]Game, errorUtils.EntityError) {
	var games []Game
	g.db.Find(&games)
	return games, nil
}

