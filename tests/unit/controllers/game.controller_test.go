package controllers

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/router"
	"GamesAPI/src/services"
	"GamesAPI/src/utils"
	"GamesAPI/src/utils/errorUtils"
	"GamesAPI/tests/unit/mocks"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GameControllerTestSuite struct {
	suite.Suite
	mockService mocks.GameServiceMockInterface
	r           *gin.Engine
	rr          *httptest.ResponseRecorder
}

func TestGamesControllerTestSuite(t *testing.T) {
	suite.Run(t, new(GameControllerTestSuite))
}

func (s *GameControllerTestSuite) SetupSuite() {
	mock := &mocks.GameServiceMock{}
	s.mockService = mock
	services.GamesService = mock
	s.r = gin.Default()
	router.InitAllGameRoutes(s.r.Group(""))
}

func (s *GameControllerTestSuite) BeforeTest(_, _ string) {
	s.rr = httptest.NewRecorder()
}

func (s *GameControllerTestSuite) TestGetGame_Success() {
	s.mockService.SetGetGame(func(id uint64) (*domain.Game, errorUtils.EntityError) {
		return &domain.Game{
			ID:          1,
			Title:       "Rocket League",
			Developer:   "Psyonix",
			Publisher:   "Psyonix",
			ReleaseDate: utils.GetDate("2015-07-07"),
		}, nil
	})
	gameIdParam := "1"
	req, _ := http.NewRequest(http.MethodGet, "/games/"+gameIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	var game domain.Game
	err := json.Unmarshal(s.rr.Body.Bytes(), &game)
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, game)
	assert.EqualValues(t, http.StatusOK, s.rr.Code)
	assert.EqualValues(t, 1, game.ID)
	assert.EqualValues(t, "Rocket League", game.Title)
	assert.EqualValues(t, "Psyonix", game.Developer)
	assert.EqualValues(t, "Psyonix", game.Publisher)
}

func (s *GameControllerTestSuite) TestGetGame_InvalidId() {
	gameIdParam := "abc"
	req, _ := http.NewRequest(http.MethodGet, "/games/"+gameIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "game id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func (s *GameControllerTestSuite) TestGetGame_NotFound() {
	s.mockService.SetGetGame(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return nil, errorUtils.NewNotFoundError("game not found")
	})
	gameIdParam := "1"
	req, _ := http.NewRequest(http.MethodGet, "/games/"+gameIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusNotFound, apiErr.Status())
	assert.EqualValues(t, "game not found", apiErr.Message())
	assert.EqualValues(t, "not_found", apiErr.Error())
}

func (s *GameControllerTestSuite) TestGetGame_DatabaseError() {
	s.mockService.SetGetGame(func(u uint64) (*domain.Game, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("database error")
	})
	gameIdParam := "1"
	req, _ := http.NewRequest(http.MethodGet, "/games/"+gameIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "database error", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

func (s *GameControllerTestSuite) TestCreateGame_Success() {
	s.mockService.SetCreateGame(func(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
		return &domain.Game{
			ID:          1,
			Title:       "Rocket League",
			Developer:   "Psyonix",
			Publisher:   "Psyonix",
			ReleaseDate: utils.GetDate("2015-07-07"),
		}, nil
	})
	jsonBody := `{"title":"Rocket League", "developer":"Psyonix", "publisher":"Psyonix"}`
	req, err := http.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	var game domain.Game
	err = json.Unmarshal(s.rr.Body.Bytes(), &game)
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, game)
	assert.EqualValues(t, http.StatusCreated, s.rr.Code)
	assert.EqualValues(t, 1, game.ID)
	assert.EqualValues(t, "Rocket League", game.Title)
	assert.EqualValues(t, "Psyonix", game.Developer)
	assert.EqualValues(t, "Psyonix", game.Publisher)
}

func (s *GameControllerTestSuite) TestCreateGame_InvalidJsonBadFieldType() {
	//here we put a number instead of string for title. we expect an invalid json error
	jsonBody := `{"title":123456, "developer":"Psyonix", "publisher":"Psyonix"}`
	req, err := http.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)
	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func (s *GameControllerTestSuite) TestCreateGame_InvalidJsonMissingField() {
	//here we put a typo in 'title' field. we expect an invalid json error.
	jsonBody := `{"titl":"Rocket League", "developer":"Psyonix", "publisher":"Psyonix"}`
	req, err := http.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)
	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func (s *GameControllerTestSuite) TestUpdateGame_Success() {
	s.mockService.SetUpdateGame(func(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
		return &domain.Game{
			ID:          1,
			Title:       "Rocket League",
			Developer:   "Psyonix",
			Publisher:   "Psyonix",
			ReleaseDate: utils.GetDate("2015-07-07"),
		}, nil
	})
	id := "1"
	jsonBody := `{"title":"Rocket League 2", "developer":"Psyonix, but better", "publisher":"Not Psyonix"}`
	req, err := http.NewRequest(http.MethodPatch, "/games/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	var game domain.Game
	err = json.Unmarshal(s.rr.Body.Bytes(), &game)
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, game)
	assert.EqualValues(t, http.StatusOK, s.rr.Code)
	assert.EqualValues(t, 1, game.ID)
	assert.EqualValues(t, "Rocket League", game.Title)
	assert.EqualValues(t, "Psyonix", game.Developer)
	assert.EqualValues(t, "Psyonix", game.Publisher)
}

func (s *GameControllerTestSuite) TestUpdateGame_InvalidId() {
	gameIdParam := "abc"
	req, _ := http.NewRequest(http.MethodPatch, "/games/"+gameIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "game id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func (s *GameControllerTestSuite) TestUpdateGame_InvalidJsonFieldMissing() {
	//here we puroposely put a typo in 'title' field name. we expect it to return us an invalid JSON error
	jsonBody := `{"titl":"Rocket League", "developer":"Psyonix", "publisher":"Psyonix"}`
	id := "1"
	req, err := http.NewRequest(http.MethodPatch, "/games/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)
	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func (s *GameControllerTestSuite) TestUpdateGame_InvalidJsonBadFieldType() {
	//here we put a number instead of a string for the title. we expect an invalid json error
	jsonBody := `{"title":123456, "developer":"Psyonix", "publisher":"Psyonix"}`
	id := "1"
	req, err := http.NewRequest(http.MethodPatch, "/games/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)
	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func (s *GameControllerTestSuite) TestUpdateGame_ErrorUpdating() {
	s.mockService.SetUpdateGame(func(game *domain.Game) (*domain.Game, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error updating game")
	})

	id := "1"
	jsonBody := `{"title":"Rocket League 2", "developer":"Psyonix, but better", "publisher":"Not Psyonix"}`
	req, err := http.NewRequest(http.MethodPatch, "/games/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error updating game", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}

func (s *GameControllerTestSuite) TestDeleteGame_Success() {
	s.mockService.SetDelete(func(u uint64) errorUtils.EntityError {
		return nil
	})
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/games/"+id, nil)
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	var response = make(map[string]string)
	theErr := json.Unmarshal(s.rr.Body.Bytes(), &response)
	if theErr != nil {
		s.T().Errorf("could not unmarshal response: %v\n", theErr)
	}

	assert.EqualValues(s.T(), http.StatusOK, s.rr.Code)
	assert.EqualValues(s.T(), "deleted", response["status"])
}

func (s *GameControllerTestSuite) TestDeleteGame_InvalidId() {
	gameIdParam := "abc"
	req, _ := http.NewRequest(http.MethodDelete, "/games/"+gameIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "game id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func (s *GameControllerTestSuite) TestDeleteGame_Failure() {
	s.mockService.SetDelete(func(u uint64) errorUtils.EntityError {
		return errorUtils.NewInternalServerError("error deleting game")
	})
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/games/"+id, nil)
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error deleting game", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}

func (s *GameControllerTestSuite) TestGetAllGames_Success() {
	s.mockService.SetGetAll(func() ([]domain.Game, errorUtils.EntityError) {
		return []domain.Game{
			{
				ID:          1,
				Title:       "Rocket League",
				Developer:   "Psyonix",
				Publisher:   "Psyonix",
				ReleaseDate: utils.GetDate("2015-07-07"),
			},
			{
				ID:          2,
				Title:       "The Witcher 3: Wild Hunt",
				Developer:   "CD PROJEKT RED",
				Publisher:   "CD PROJEKT RED",
				ReleaseDate: utils.GetDate("2015-05-18"),
			},
		}, nil
	})

	req, err := http.NewRequest(http.MethodGet, "/games", nil)
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	var games []domain.Game
	theErr := json.Unmarshal(s.rr.Body.Bytes(), &games)
	if theErr != nil {
		s.T().Errorf("could not unmarshal response: %v\n", theErr)
	}
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, games)
	assert.EqualValues(t, 1, games[0].ID)
	assert.EqualValues(t, "Rocket League", games[0].Title)
	assert.EqualValues(t, "Psyonix", games[0].Developer)
	assert.EqualValues(t, "Psyonix", games[0].Publisher)
	assert.EqualValues(t, 2, games[1].ID)
	assert.EqualValues(t, "The Witcher 3: Wild Hunt", games[1].Title)
	assert.EqualValues(t, "CD PROJEKT RED", games[1].Developer)
	assert.EqualValues(t, "CD PROJEKT RED", games[1].Publisher)
}

func (s *GameControllerTestSuite) TestGetAllGames_Failure() {
	s.mockService.SetGetAll(func() ([]domain.Game, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error getting games")
	})
	req, err := http.NewRequest(http.MethodGet, "/games", nil)
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error getting games", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}
