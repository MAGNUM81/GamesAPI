package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils"
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type GameRepoMockInterface interface {
	SetGetGameDomain(func(uint64) (*domain.Game, errorUtils.EntityError))
	SetCreateGameDomain(func(game *domain.Game) (*domain.Game, errorUtils.EntityError))
	SetUpdateGameDomain(func(game *domain.Game) (*domain.Game, errorUtils.EntityError))
	SetDeleteGameDomain(func(id uint64) errorUtils.EntityError)
	SetGetAllGameDomain(func() ([]domain.Game, errorUtils.EntityError))
}

type gameRepoMock struct {
	getGameDomain func(id uint64) (*domain.Game, errorUtils.EntityError)
	createGameDomain func(game *domain.Game) (*domain.Game, errorUtils.EntityError)
	updateGameDomain func(game *domain.Game) (*domain.Game, errorUtils.EntityError)
	deleteGameDomain func(id uint64) errorUtils.EntityError
	getAllGamesDomain func() ([]domain.Game, errorUtils.EntityError)
}

//GameRepoMockInterface implementation, so we can swap the methods around and get the desired behavior from the repository
func (m *gameRepoMock) SetGetGameDomain(f func(uint64) (*domain.Game, errorUtils.EntityError)) {
	m.getGameDomain = f
}

func (m *gameRepoMock) SetCreateGameDomain(f func(game *domain.Game) (*domain.Game, errorUtils.EntityError)) {
	m.createGameDomain = f
}

func (m *gameRepoMock) SetUpdateGameDomain(f func(game *domain.Game) (*domain.Game, errorUtils.EntityError)) {
	m.updateGameDomain = f
}

func (m *gameRepoMock) SetDeleteGameDomain(f func(id uint64) errorUtils.EntityError) {
	m.deleteGameDomain = f
}

func (m *gameRepoMock) SetGetAllGameDomain(f func() ([]domain.Game, errorUtils.EntityError)) {
	m.getAllGamesDomain = f
}

//GameRepoInterface implementation (redirects all calls to the swappable methods)
func (m *gameRepoMock) Get(id uint64) (*domain.Game, errorUtils.EntityError){
	return m.getGameDomain(id)
}
func (m *gameRepoMock) Create(msg *domain.Game) (*domain.Game, errorUtils.EntityError){
	return m.createGameDomain(msg)
}
func (m *gameRepoMock) Update(msg *domain.Game) (*domain.Game, errorUtils.EntityError){
	return m.updateGameDomain(msg)
}
func (m *gameRepoMock) Delete(id uint64) errorUtils.EntityError {
	return m.deleteGameDomain(id)
}
func (m *gameRepoMock) GetAll() ([]domain.Game, errorUtils.EntityError) {
	return m.getAllGamesDomain()
}
func (m *gameRepoMock) Initialize(_ *gorm.DB) {}

type GameServiceTestSuite struct {
	suite.Suite
	mockRepository GameRepoMockInterface
}

func TestGameServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GameServiceTestSuite))
}

func (s *GameServiceTestSuite) SetupSuite() {
	mock := &gameRepoMock{}

	s.mockRepository = mock //set this so we can swap the methods
	domain.GameRepo = mock  //set this so the tested code calls the swapped methods
}

func (s *GameServiceTestSuite) TestGamesService_GetGame_Success() {
	s.mockRepository.SetGetGameDomain(func(gameId uint64) (*domain.Game, errorUtils.EntityError) {
		return &domain.Game{
			ID:          1,
			Title:       "Rocket League",
			Developer:   "Psyonix",
			Publisher:   "Psyonix",
			ReleaseDate: utils.GetDate("2015-07-07"),
			CreatedAt: tm,
		}, nil
	})
	game, err := services.GamesService.GetGame(1)
	t := s.T()
	assert.NotNil(t, game)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, game.ID)
	assert.EqualValues(t, "Rocket League", game.Title)
	assert.EqualValues(t, "Psyonix", game.Developer)
	assert.EqualValues(t, "Psyonix", game.Publisher)
	assert.EqualValues(t, tm, game.CreatedAt)
}

func (s *GameServiceTestSuite) TestGamesService_GetGame_NotFound() {
	expectedError := errorUtils.NewNotFoundError("game was not found")
	s.mockRepository.SetGetGameDomain(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return nil, expectedError
	})
	game, err := services.GamesService.GetGame(1)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), game)
	assert.Equal(s.T(), expectedError, err)
}

func (s *GameServiceTestSuite) TestGamesService_CreateGame_Success() {
	expectedGame := &domain.Game{
		ID:          1,
		Title:       "Rocket League",
		Developer:   "Psyonix",
		Publisher:   "Psyonix",
		ReleaseDate: utils.GetDate("2015-07-07"),
		CreatedAt: tm,
	}
	s.mockRepository.SetCreateGameDomain(func(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
		return expectedGame, nil
	})
	request := &domain.Game{
		Title:       "Rocket League",
		Developer:   "Psyonix",
		Publisher:   "Psyonix",
		ReleaseDate: utils.GetDate("2015-07-07"),
		CreatedAt: tm,
	}
	game, err := services.GamesService.CreateGame(request)
	assert.NotNil(s.T(), game)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedGame, game)
}

func (s *GameServiceTestSuite) TestGamesService_CreateGame_InvalidRequest() {
	tests := []struct {
		request *domain.Game
		expectedError errorUtils.EntityError
	}{
		{
			request: &domain.Game{
				Title:       "",
				Developer:   "Psyonix",
				Publisher:   "Psyonix",
				ReleaseDate: utils.GetDate("2015-07-07"),
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("Game title cannot be empty"),
		},
		{
			request: &domain.Game{
				Title:       "Rocket League",
				Developer:   "",
				Publisher:   "Psyonix",
				ReleaseDate: utils.GetDate("2015-07-07"),
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("Game developer cannot be empty"),
		},
		{
			request: &domain.Game{
				Title:       "Rocket League",
				Developer:   "Psyonix",
				Publisher:   "",
				ReleaseDate: utils.GetDate("2015-07-07"),
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("Game publisher cannot be empty"),
		},
	}
	for _, tt := range tests {
		msg, err := services.GamesService.CreateGame(tt.request)
		assert.Nil(s.T(), msg)
		assert.NotNil(s.T(), err)
		assert.Equal(s.T(), tt.expectedError, err)
	}
}

func (s *GameServiceTestSuite) TestGamesService_UpdateGame_Success() {
	before := &domain.Game{
		ID:          1,
		Title:       "Rocket League",
		Developer:   "Psyonix",
		Publisher:   "Psyonix",
		ReleaseDate: utils.GetDate("2015-07-07"),
		CreatedAt: tm,
	}
	expectedAfter := &domain.Game{
		ID:          1,
		Title:       "Rocket League After",
		Developer:   "Psyonix After",
		Publisher:   "Psyonix After",
		ReleaseDate: utils.GetDate("2051-07-07"),
		CreatedAt: tm,
	}
	s.mockRepository.SetGetGameDomain(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return before, nil
	})
	s.mockRepository.SetUpdateGameDomain(func(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
		return expectedAfter, nil
	})

	request := expectedAfter

	game, err := services.GamesService.UpdateGame(request)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), game)
	assert.Equal(s.T(), expectedAfter, game)
}

func (s *GameServiceTestSuite) TestGamesService_UpdateGame_InvalidRequest() {
	tests := []struct {
		request *domain.Game
		expectedError errorUtils.EntityError
	}{
		{
			request: &domain.Game{
				ID: 		 1,
				Title:       "",
				Developer:   "Psyonix",
				Publisher:   "Psyonix",
				ReleaseDate: utils.GetDate("2015-07-07"),
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("Game title cannot be empty"),
		},
		{
			request: &domain.Game{
				ID: 		 1,
				Title:       "Rocket League",
				Developer:   "",
				Publisher:   "Psyonix",
				ReleaseDate: utils.GetDate("2015-07-07"),
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("Game developer cannot be empty"),
		},
		{
			request: &domain.Game{
				ID: 		 1,
				Title:       "Rocket League",
				Developer:   "Psyonix",
				Publisher:   "",
				ReleaseDate: utils.GetDate("2015-07-07"),
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("Game publisher cannot be empty"),
		},
	}
	for _, tt := range tests {
		msg, err := services.GamesService.UpdateGame(tt.request)
		assert.Nil(s.T(), msg)
		assert.NotNil(s.T(), err)
		assert.Equal(s.T(), tt.expectedError, err)
	}
}

func (s *GameServiceTestSuite) TestGamesService_UpdateGame_FailureGettingFormerGame() {
	s.mockRepository.SetGetGameDomain(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error getting game")
	})
	request := &domain.Game{
		Title:       "Rocket League",
		Developer:   "Psyonix",
		Publisher:   "Psyonix",
	}
	msg, err := services.GamesService.UpdateGame(request)
	t := s.T()
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error getting game", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func (s *GameServiceTestSuite) TestGamesService_UpdateGame_FailureUpdatingGame() {
	s.mockRepository.SetGetGameDomain(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return &domain.Game{
			ID:    1,
			Title:       "Rocket League",
			Developer:   "Psyonix",
			Publisher:   "Psyonix",
		}, nil
	})
	s.mockRepository.SetUpdateGameDomain(func(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error updating game")
	})

	request := &domain.Game{
		ID:    1,
		Title:       "Rocket League AAA",
		Developer:   "Psyonix AAA",
		Publisher:   "Psyonix AAA",
	}
	msg, err := services.GamesService.UpdateGame(request)
	t:= s.T()
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error updating game", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func (s *GameServiceTestSuite) TestGamesService_DeleteGame_Success() {
	s.mockRepository.SetGetGameDomain(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return &domain.Game{
			ID:    1,
			Title:       "Rocket League",
			Developer:   "Psyonix",
			Publisher:   "Psyonix",
		}, nil
	})
	s.mockRepository.SetDeleteGameDomain(func(_ uint64) errorUtils.EntityError {
		return nil
	})

	err := services.GamesService.DeleteGame(1)
	assert.Nil(s.T(), err)
}

func (s *GameServiceTestSuite) TestGamesService_DeleteGame_ErrorGettingGame() {
	expectedError := errorUtils.NewInternalServerError("Something went wrong fetching game")
	s.mockRepository.SetGetGameDomain(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return nil, expectedError
	})
	err := services.GamesService.DeleteGame(1)
	t := s.T()
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func (s *GameServiceTestSuite) TestGamesService_DeleteGame_ErrorDeletingGame(){
	expectedError := errorUtils.NewInternalServerError("error deleting message")
	s.mockRepository.SetGetGameDomain(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return &domain.Game{
			ID:    1,
			Title:       "Rocket League",
			Developer:   "Psyonix",
			Publisher:   "Psyonix",
		}, nil
	})
	s.mockRepository.SetDeleteGameDomain(func(id uint64) errorUtils.EntityError {
		return expectedError
	})

	err := services.GamesService.DeleteGame(1)
	t := s.T()
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func (s *GameServiceTestSuite) TestGamesService_GetAll_Success() {
	s.mockRepository.SetGetAllGameDomain(func() ([]domain.Game, errorUtils.EntityError) {
		return []domain.Game{
			{
				ID: 1,
				Title:       "Rocket League1",
				Developer:   "Psyonix1",
				Publisher:   "Psyonix1",
			},
			{
				ID: 2,
				Title:       "Rocket League2",
				Developer:   "Psyonix2",
				Publisher:   "Psyonix2",
			},
			{
				ID: 3,
				Title:       "Rocket League3",
				Developer:   "Psyonix3",
				Publisher:   "Psyonix3",
			},
		}, nil
	})
	games, err := services.GamesService.GetAllGames()
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, games)
	assert.EqualValues(t, games[0].ID, 1)
	assert.EqualValues(t, games[0].Title, "Rocket League1")
	assert.EqualValues(t, games[0].Developer, "Psyonix1")
	assert.EqualValues(t, games[0].Publisher, "Psyonix1")
	assert.EqualValues(t, games[1].ID, 2)
	assert.EqualValues(t, games[1].Title, "Rocket League2")
	assert.EqualValues(t, games[1].Developer, "Psyonix2")
	assert.EqualValues(t, games[1].Publisher, "Psyonix2")
	assert.EqualValues(t, games[2].ID, 3)
	assert.EqualValues(t, games[2].Title, "Rocket League3")
	assert.EqualValues(t, games[2].Developer, "Psyonix3")
	assert.EqualValues(t, games[2].Publisher, "Psyonix3")
}

func (s *GameServiceTestSuite) TestGamesService_GetAllGames_ErrorGettingGames() {
	expectedErr := errorUtils.NewInternalServerError("error getting games")
	s.mockRepository.SetGetAllGameDomain(func() ([]domain.Game, errorUtils.EntityError) {
		return nil, expectedErr
	})

	games, err := services.GamesService.GetAllGames()
	t := s.T()
	assert.NotNil(t, err)
	assert.Nil(t, games)
	assert.Equal(t, expectedErr, err)
}
