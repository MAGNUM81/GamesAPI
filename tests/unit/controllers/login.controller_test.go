package controllers

import (
	"GamesAPI/src/controllers"
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"GamesAPI/tests/unit/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type LoginControllerTestSuite struct {
	suite.Suite
	mockUsersService          mocks.UserServiceMockInterface
	mockUserSessionService    mocks.UserSessionServiceMockInterface
	mockAuthenticationService mocks.AuthenticationServiceMockInterface
	r                         *gin.Engine
	rr                        *httptest.ResponseRecorder
}

func TestLoginControllerTestSuite(t *testing.T) {
	suite.Run(t, new(LoginControllerTestSuite))
}

func (s *LoginControllerTestSuite) SetupSuite() {
	usersMock := &mocks.UserServiceMock{}
	userSeshMock := &mocks.UserSessionServiceMock{}
	authMock := &mocks.AuthenticationServiceMock{}
	s.mockUserSessionService = userSeshMock
	s.mockAuthenticationService = authMock
	s.mockUsersService = usersMock
	services.UserSessionService = userSeshMock
	services.AuthenticationService = authMock
	services.UsersService = usersMock
	s.r = gin.Default()
	s.r.GET("/auth/login", controllers.LoginController)
}

func (s *LoginControllerTestSuite) BeforeTest(_, _ string) {
	s.rr = httptest.NewRecorder()
}

func (s *LoginControllerTestSuite) TestLogin_NoHeader() {
	req, _ := http.NewRequest(http.MethodGet, "/auth/login", nil)
	s.r.ServeHTTP(s.rr, req)
	t := s.T()
	assert.Equal(t, http.StatusBadRequest, s.rr.Code)
}

func (s *LoginControllerTestSuite) TestLogin_BadHeader() {
	req, _ := http.NewRequest(http.MethodGet, "/auth/login", nil)
	req.Header.Add("Authorization", "username;password;favorite_fantasy_animal")
	s.r.ServeHTTP(s.rr, req)
	t := s.T()
	assert.Equal(t, http.StatusBadRequest, s.rr.Code)
}

func (s *LoginControllerTestSuite) TestLogin_ErrGettingAllUsers() {
	s.mockUsersService.SetGetAll(func() ([]domain.User, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("Sumtin wong")
	})

	req, _ := http.NewRequest(http.MethodGet, "/auth/login", nil)
	req.Header.Add("Authorization", "username;password")
	s.r.ServeHTTP(s.rr, req)
	t := s.T()
	assert.Equal(t, http.StatusInternalServerError, s.rr.Code)
}

func (s *LoginControllerTestSuite) TestLogin_NoUserMatchesEmail() {
	s.mockUsersService.SetGetAll(func() ([]domain.User, errorUtils.EntityError) {
		return []domain.User{
			{
				ID:           1,
				Name:         "dev",
				Email:        "dev@golang.com",
				PasswordHash: "epsteindidntkillhimself",
			},
		}, nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/auth/login", nil)
	req.Header.Add("Authorization", "golang@dev.com;some_password")
	s.r.ServeHTTP(s.rr, req)
	t := s.T()
	assert.Equal(t, http.StatusUnauthorized, s.rr.Code)
}

func (s *LoginControllerTestSuite) TestLogin_PasswordInvalid() {
	s.mockUsersService.SetGetAll(func() ([]domain.User, errorUtils.EntityError) {
		return []domain.User{
			{
				ID:           1,
				Name:         "dev",
				Email:        "dev@golang.com",
				PasswordHash: "epsteindidntkillhimself",
			},
		}, nil
	})

	s.mockAuthenticationService.SetValidatePassword(func(plainPassword []byte, hashedPassword string) (bool, error) {
		return false, nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/auth/login", nil)
	req.Header.Add("Authorization", "dev@golang.com;some_password")
	s.r.ServeHTTP(s.rr, req)
	t := s.T()
	assert.Equal(t, http.StatusUnauthorized, s.rr.Code)
}

func (s *LoginControllerTestSuite) TestLogin_CouldNotGenerateSessionToken() {
	s.mockUsersService.SetGetAll(func() ([]domain.User, errorUtils.EntityError) {
		return []domain.User{
			{
				ID:           1,
				Name:         "dev",
				Email:        "dev@golang.com",
				PasswordHash: "epsteindidntkillhimself",
			},
		}, nil
	})

	s.mockAuthenticationService.SetValidatePassword(func(plainPassword []byte, hashedPassword string) (bool, error) {
		return true, nil
	})

	s.mockUserSessionService.SetGenerateSessionToken(func(userId uint64, expireAt time.Time) (string, error) {
		return "", errors.New("could not generate token")
	})

	req, _ := http.NewRequest(http.MethodGet, "/auth/login", nil)
	req.Header.Add("Authorization", "dev@golang.com;some_password")
	s.r.ServeHTTP(s.rr, req)
	t := s.T()
	assert.Equal(t, http.StatusInternalServerError, s.rr.Code)
}

func (s *LoginControllerTestSuite) TestLogin_CouldNotCreateUserSession() {
	s.mockUsersService.SetGetAll(func() ([]domain.User, errorUtils.EntityError) {
		return []domain.User{
			{
				ID:           1,
				Name:         "dev",
				Email:        "dev@golang.com",
				PasswordHash: "epsteindidntkillhimself",
			},
		}, nil
	})

	s.mockAuthenticationService.SetValidatePassword(func(plainPassword []byte, hashedPassword string) (bool, error) {
		return true, nil
	})

	s.mockUserSessionService.SetGenerateSessionToken(func(userId uint64, expireAt time.Time) (string, error) {
		return "some_token", nil
	})

	s.mockUserSessionService.SetCreateSession(func(token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("could not create user session")
	})

	req, _ := http.NewRequest(http.MethodGet, "/auth/login", nil)
	req.Header.Add("Authorization", "dev@golang.com;some_password")
	s.r.ServeHTTP(s.rr, req)
	t := s.T()
	assert.Equal(t, http.StatusInternalServerError, s.rr.Code)
}

func (s *LoginControllerTestSuite) TestLogin_Success() {
	s.mockUsersService.SetGetAll(func() ([]domain.User, errorUtils.EntityError) {
		return []domain.User{
			{
				ID:           1,
				Name:         "dev",
				Email:        "dev@golang.com",
				PasswordHash: "epsteindidntkillhimself",
			},
		}, nil
	})

	s.mockAuthenticationService.SetValidatePassword(func(plainPassword []byte, hashedPassword string) (bool, error) {
		return true, nil
	})

	s.mockUserSessionService.SetGenerateSessionToken(func(userId uint64, expireAt time.Time) (string, error) {
		return "some_token", nil
	})

	s.mockUserSessionService.SetCreateSession(func(token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError) {
		return &domain.UserSession{
			Token:     "some_token",
			UserId:    1,
			ExpiresAt: time.Now().Add(time.Minute * 10).UnixNano(),
		}, nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/auth/login", nil)
	req.Header.Add("Authorization", "dev@golang.com;some_password")
	s.r.ServeHTTP(s.rr, req)
	t := s.T()
	assert.Equal(t, http.StatusOK, s.rr.Code)
	//check if we got the token back in the response
	assert.Equal(t, "some_token", s.rr.Header().Get("Authorization"))
}
