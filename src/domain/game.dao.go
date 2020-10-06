package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
)

var (
	GameRepo GameRepoInterface = &gameRepo{}
)

type GameRepoInterface interface {
	Get(uint64) (*Game, errorUtils.GameError)
	Create(*Game) (*Game, errorUtils.GameError)
	Update(*Game) (*Game, errorUtils.GameError)
	Delete(uint64) errorUtils.GameError
	GetAll() ([]Game, errorUtils.GameError)
	Initialize(*gorm.DB)
}

type gameRepo struct {
	db *gorm.DB
}

func (g *gameRepo) Initialize(db *gorm.DB) {
	db.AutoMigrate(&Game{})
}

func NewGameRepository(db *gorm.DB) GameRepoInterface {
	return &gameRepo{db: db}
}

func (g *gameRepo) Get(gameId uint64) (*Game, errorUtils.GameError) {
	var game Game
	if err := g.db.Where("id = ?", gameId).First(&game).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	return &game, nil
}

func (g *gameRepo) Create(game *Game) (*Game, errorUtils.GameError) {
	g.db.Create(game)
	return game, nil
}

func (g *gameRepo) Update(game *Game) (*Game, errorUtils.GameError) {
	if err := g.db.Where("id = ?", game.ID).First(&game).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	g.db.Save(*game)
	return game, nil
}

func (g *gameRepo) Delete(gameId uint64) errorUtils.GameError {
	var game Game
	if err := g.db.Where("id = ?", gameId).First(&game).Error; err != nil {
		return errorUtils.NewNotFoundError(err.Error())
	}
	g.db.Delete(&game)
	return nil
}

func (g *gameRepo) GetAll() ([]Game, errorUtils.GameError) {
	var games []Game
	g.db.Find(&games)
	return games, nil
}

