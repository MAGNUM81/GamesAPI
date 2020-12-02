package middleware

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/middleware"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"GamesAPI/tests/unit/mocks"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type UserSessionHandlerTestSuite struct {
	suite.Suite
	mockService mocks.UserSessionServiceMockInterface
	r           *gin.Engine
	rr          *httptest.ResponseRecorder
}

func UserSessionBidonHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Message": "Success"})
}

func (s *UserSessionHandlerTestSuite) BeforeTest(_, _ string) {
	s.rr = httptest.NewRecorder()
}

func TestUserSessionHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserSessionHandlerTestSuite))
}

func (s *UserSessionHandlerTestSuite) SetupSuite() {
	mock := &mocks.UserSessionServiceMock{}

	s.mockService = mock
	services.UserSessionService = mock
	s.r = gin.Default()
	s.r.Use(middleware.UserSessionHandler)
	s.r.GET("/", UserSessionBidonHandler)
}

func (s *UserSessionHandlerTestSuite) TestUserSessionHandler_NoAuthHeader() {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s.r.ServeHTTP(s.rr, req)
	assert.Equal(s.T(), http.StatusBadRequest, s.rr.Code)
}

func (s *UserSessionHandlerTestSuite) TestUserSessionHandler_BadAuthHeader1() {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "This is a bad Auth Header")
	s.r.ServeHTTP(s.rr, req)
	assert.Equal(s.T(), http.StatusBadRequest, s.rr.Code)
}

func (s *UserSessionHandlerTestSuite) TestUserSessionHandler_BadAuthHeader2() {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer ")
	s.r.ServeHTTP(s.rr, req)
	assert.Equal(s.T(), http.StatusBadRequest, s.rr.Code)
}

func (s *UserSessionHandlerTestSuite) TestUserSessionHandler_SessionNotExists() {
	s.mockService.SetExistsSession(func(key string) bool {
		return false
	})
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer 123456")
	s.r.ServeHTTP(s.rr, req)
	assert.Equal(s.T(), http.StatusUnauthorized, s.rr.Code)
	responseWWWAuth := s.rr.Header().Get("Www-ValidatePassword")
	assert.True(s.T(), responseWWWAuth != "")
}

func (s *UserSessionHandlerTestSuite) TestUserSessionHandler_IsSessionExpired_Failure() {
	s.mockService.SetExistsSession(func(key string) bool {
		return true
	})
	s.mockService.SetIsSessionExpired(func(key string, currentTime time.Time) (bool, errorUtils.EntityError) {
		return true, errorUtils.NewNotFoundError(fmt.Sprintf("token with key %s does not exist", key))
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer 123456")
	s.r.ServeHTTP(s.rr, req)
	assert.Equal(s.T(), http.StatusNotFound, s.rr.Code)
}

func (s *UserSessionHandlerTestSuite) TestUserSessionHandler_IsSessionExpired_Success() {
	s.mockService.SetExistsSession(func(key string) bool {
		return true
	})
	s.mockService.SetIsSessionExpired(func(key string, currentTime time.Time) (bool, errorUtils.EntityError) {
		return true, nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer 123456")
	s.r.ServeHTTP(s.rr, req)
	assert.Equal(s.T(), http.StatusUnauthorized, s.rr.Code)
	responseWWWAuth := s.rr.Header().Get("Www-ValidatePassword")
	assert.True(s.T(), responseWWWAuth != "")
}

func (s *UserSessionHandlerTestSuite) TestUserSessionHandler_SessionIsNotExpired() {
	s.mockService.SetExistsSession(func(key string) bool {
		return true
	})
	s.mockService.SetIsSessionExpired(func(key string, currentTime time.Time) (bool, errorUtils.EntityError) {
		return false, nil
	})

	s.mockService.SetGetSession(func(key string) (*domain.UserSession, errorUtils.EntityError) {
		return &domain.UserSession{
			Token:     "12345",
			UserId:    1,
			ExpiresAt: time.Now().Add(time.Minute * 1).UnixNano(),
		}, nil
	})

	req, _ :=http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer 123456")
	s.r.ServeHTTP(s.rr, req)
	assert.Equal(s.T(), http.StatusOK, s.rr.Code)
}
