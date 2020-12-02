package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
)

type GameRepoMockInterface interface {
	SetGetGameDomain(func(uint64) (*domain.Game, errorUtils.EntityError))
	SetCreateGameDomain(func(game *domain.Game) (*domain.Game, errorUtils.EntityError))
	SetUpdateGameDomain(func(game *domain.Game) (*domain.Game, errorUtils.EntityError))
	SetDeleteGameDomain(func(id uint64) errorUtils.EntityError)
	SetGetAllGameDomain(func() ([]domain.Game, errorUtils.EntityError))
}

type GameRepoMock struct {
	getGameDomain     func(id uint64) (*domain.Game, errorUtils.EntityError)
	createGameDomain  func(game *domain.Game) (*domain.Game, errorUtils.EntityError)
	updateGameDomain  func(game *domain.Game) (*domain.Game, errorUtils.EntityError)
	deleteGameDomain  func(id uint64) errorUtils.EntityError
	getAllGamesDomain func() ([]domain.Game, errorUtils.EntityError)
}

//GameRepoMockInterface implementation, so we can swap the methods around and get the desired behavior from the repository
func (m *GameRepoMock) SetGetGameDomain(f func(uint64) (*domain.Game, errorUtils.EntityError)) {
	m.getGameDomain = f
}

func (m *GameRepoMock) SetCreateGameDomain(f func(game *domain.Game) (*domain.Game, errorUtils.EntityError)) {
	m.createGameDomain = f
}

func (m *GameRepoMock) SetUpdateGameDomain(f func(game *domain.Game) (*domain.Game, errorUtils.EntityError)) {
	m.updateGameDomain = f
}

func (m *GameRepoMock) SetDeleteGameDomain(f func(id uint64) errorUtils.EntityError) {
	m.deleteGameDomain = f
}

func (m *GameRepoMock) SetGetAllGameDomain(f func() ([]domain.Game, errorUtils.EntityError)) {
	m.getAllGamesDomain = f
}

//GameRepoInterface implementation (redirects all calls to the swappable methods)
func (m *GameRepoMock) Get(id uint64) (*domain.Game, errorUtils.EntityError) {
	return m.getGameDomain(id)
}
func (m *GameRepoMock) Create(msg *domain.Game) (*domain.Game, errorUtils.EntityError) {
	return m.createGameDomain(msg)
}
func (m *GameRepoMock) Update(msg *domain.Game) (*domain.Game, errorUtils.EntityError) {
	return m.updateGameDomain(msg)
}
func (m *GameRepoMock) Delete(id uint64) errorUtils.EntityError {
	return m.deleteGameDomain(id)
}
func (m *GameRepoMock) GetAll() ([]domain.Game, errorUtils.EntityError) {
	return m.getAllGamesDomain()
}
func (m *GameRepoMock) Initialize(_ *gorm.DB) {}
