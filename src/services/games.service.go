package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
)

var (
	GamesService gamesServiceInterface = &gamesService{}
)

type gamesService struct{}

type gamesServiceInterface interface {
	GetGame(uint64) (*domain.Game, errorUtils.EntityError)
	CreateGame(*domain.Game) (*domain.Game, errorUtils.EntityError)
	UpdateGame(game *domain.Game) (*domain.Game, errorUtils.EntityError)
	DeleteGame(uint64) errorUtils.EntityError
	GetAllGames() ([]domain.Game, errorUtils.EntityError)
}

func (g *gamesService) GetGame(gameId uint64) (*domain.Game, errorUtils.EntityError) {
	game, err := domain.GameRepo.Get(gameId)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g *gamesService) GetAllGames() ([]domain.Game, errorUtils.EntityError) {
	games, err := domain.GameRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (g *gamesService) CreateGame(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
	game, err := domain.GameRepo.Create(game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g *gamesService) UpdateGame(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
	current, err := domain.GameRepo.Get(uint64(game.ID))
	if err != nil {
		return nil, err
	}
	current.Developer = game.Developer
	current.Publisher = game.Publisher
	current.Title = game.Title
	current.ReleaseDate = game.ReleaseDate

	updatedGame, err := domain.GameRepo.Update(current)
	if err != nil {
		return nil, err
	}
	return updatedGame, nil
}

func (g *gamesService) DeleteGame(gameId uint64) errorUtils.EntityError {
	game, err := domain.GameRepo.Get(gameId)
	if err != nil {
		return err
	}

	deleteErr := domain.GameRepo.Delete(uint64(game.ID))
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}