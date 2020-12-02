package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
)

type GameServiceMockInterface interface {
	SetGetGame(func(uint64) (*domain.Game, errorUtils.EntityError))
	SetCreateGame(func(*domain.Game) (*domain.Game, errorUtils.EntityError))
	SetUpdateGame(func(*domain.Game) (*domain.Game, errorUtils.EntityError))
	SetDelete(func(uint64) errorUtils.EntityError)
	SetGetAll(func() ([]domain.Game, errorUtils.EntityError))
}

type GameServiceMock struct {
	getGameService    func(uint64) (*domain.Game, errorUtils.EntityError)
	createGameService func(*domain.Game) (*domain.Game, errorUtils.EntityError)
	updateGameService func(*domain.Game) (*domain.Game, errorUtils.EntityError)
	deleteGameService func(uint64) errorUtils.EntityError
	getAllGameService func() ([]domain.Game, errorUtils.EntityError)
}

func (u *GameServiceMock) GetGame(id uint64) (*domain.Game, errorUtils.EntityError) {
	return u.getGameService(id)
}

func (u *GameServiceMock) CreateGame(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
	return u.createGameService(game)
}

func (u *GameServiceMock) UpdateGame(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
	return u.updateGameService(game)
}

func (u *GameServiceMock) DeleteGame(id uint64) errorUtils.EntityError {
	return u.deleteGameService(id)
}

func (u *GameServiceMock) GetAllGames() ([]domain.Game, errorUtils.EntityError) {
	return u.getAllGameService()
}

func (u *GameServiceMock) SetGetGame(f func(uint64) (*domain.Game, errorUtils.EntityError)) {
	u.getGameService = f
}

func (u *GameServiceMock) SetCreateGame(f func(*domain.Game) (*domain.Game, errorUtils.EntityError)) {
	u.createGameService = f
}

func (u *GameServiceMock) SetUpdateGame(f func(*domain.Game) (*domain.Game, errorUtils.EntityError)) {
	u.updateGameService = f
}

func (u *GameServiceMock) SetDelete(f func(uint64) errorUtils.EntityError) {
	u.deleteGameService = f
}

func (u *GameServiceMock) SetGetAll(f func() ([]domain.Game, errorUtils.EntityError)) {
	u.getAllGameService = f
}
