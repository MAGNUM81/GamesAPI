package games

import (
	"fmt"
	"sync"
)

//inspired from: https://adodd.net/post/go-ddd-repository-pattern/
type GameRepo interface {
	All() []Game
	Get(id uint)
	RequestCreate(input *CreateGameInput) (Game, error)
	RequestUpdate(input *UpdateGameInput) (Game, error)
	RequestDelete(id uint) error
}

type InMemoryGameRepo struct {
	nextId uint
	byId   map[uint]Game
	s      sync.RWMutex
}

func NewGameRepo() *InMemoryGameRepo {
	return &InMemoryGameRepo{
		nextId: 0,
		byId:   make(map[uint]Game),
	}
}

func (repo *InMemoryGameRepo) All() []Game {
	repo.s.RLock()
	defer repo.s.RUnlock()

	gameList := make([]Game, 0, len(repo.byId))

	for _, value := range repo.byId {
		gameList = append(gameList, value)
	}
	return gameList
}

func (repo *InMemoryGameRepo) RequestCreate(input Game) (Game, error) {
	repo.s.Lock()
	defer repo.s.Unlock()

	game := Game{
		Id:          repo.nextId,
		Title:       input.Title,
		ReleaseDate: input.ReleaseDate,
		Developer: input.Developer,
		Publisher: input.Publisher,
	}

	repo.byId[game.Id] = game
	repo.nextId += 1

	return game, nil
}

func (repo *InMemoryGameRepo) RequestDelete(id uint) error {
	_, ok := repo.byId[id]
	if !ok {
		return fmt.Errorf("game does not exist for id: %d", id)
	}

	delete(repo.byId, id)
	return nil
}

func (repo *InMemoryGameRepo) RequestUpdate(input Game) (Game, error) {
	repo.s.Lock()
	defer repo.s.Unlock()
	game, ok := repo.byId[input.Id]
	if !ok {
		return Game{}, fmt.Errorf("game does not exist for id: %d", input.Id)
	}

	game.Publisher = input.Publisher
	game.Developer = input.Developer
	game.ReleaseDate = input.ReleaseDate
	game.Title = input.Title

	repo.byId[input.Id] = game
	return game, nil
}
