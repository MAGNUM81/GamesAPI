package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"GamesAPI/tests/unit/mocks"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type UserSessionServiceTestSuite struct {
	suite.Suite
	mockRepo mocks.UserSessionRepoMockInterface
}

func TestUserSessionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserSessionServiceTestSuite))
}

func (s *UserSessionServiceTestSuite) SetupSuite() {
	mock := &mocks.UserSessionRepoMock{}
	s.mockRepo = mock
	domain.UserSessionRepo = mock
}

var (
	testTimeNow = time.Now()
	testTime    = testTimeNow.UnixNano()
	testUserId  = uint64(42)
	testToken   = "1234"
	testSession = &domain.UserSession{Token: testToken, UserId: testUserId, ExpiresAt: testTime}
)

func (s *UserSessionServiceTestSuite) TestGenerateSessionToken_Success() {
	token, err := services.UserSessionService.GenerateSessionToken(testUserId, testTimeNow)
	assert.NotNil(s.T(), token)
	assert.Nil(s.T(), err)

}

func (s *UserSessionServiceTestSuite) TestCreateSession_Success() {
	s.mockRepo.SetExists(func(key string) bool {
		return false
	})
}

func (s *UserSessionServiceTestSuite) TestCreateSession_Failure_InvalidToken() {
	expected := errorUtils.NewUnprocessableEntityError("Token cannot be empty")
	sesh, err := services.UserSessionService.CreateSession(&domain.UserSession{
		Token:     "",
		UserId:    uint64(3),
		ExpiresAt: testTime,
	})

	assert.Nil(s.T(), sesh)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expected, err)
}

func (s *UserSessionServiceTestSuite) TestCreateSession_Failure_TokenAlreadyExists() {
	expected := errorUtils.NewUnprocessableEntityError(fmt.Sprintf("token with key %s already exists", testToken))
	s.mockRepo.SetExists(func(key string) bool {
		return true
	})

	sesh, err := services.UserSessionService.CreateSession(testSession)
	assert.Nil(s.T(), sesh)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expected, err)
}

func (s *UserSessionServiceTestSuite) TestExistsSession_Success() {
	s.mockRepo.SetExists(func(key string) bool {
		return true
	})

	actual := services.UserSessionService.ExistsSession(testToken)
	assert.Equal(s.T(), true, actual)
}

func (s *UserSessionServiceTestSuite) TestDeleteSession_Success() {
	s.mockRepo.SetExists(func(key string) bool {
		return true
	})

	s.mockRepo.SetDelete(func(key string) errorUtils.EntityError {
		return nil
	})

	err := services.UserSessionService.DeleteSession(testToken)

	assert.Nil(s.T(), err)
}

func (s *UserSessionServiceTestSuite) TestDeleteSession_Failure_NotFound() {
	expected := errorUtils.NewNotFoundError(fmt.Sprintf("token with key %s does not exist", testToken))
	s.mockRepo.SetExists(func(key string) bool {
		return false
	})

	actual := services.UserSessionService.DeleteSession(testToken)

	assert.NotNil(s.T(), actual)
	assert.Equal(s.T(), expected, actual)
}

func (s *UserSessionServiceTestSuite) TestDeleteSession_Failure_RepoError() {
	expected := errorUtils.NewInternalServerError("error when deleting session from repo")
	s.mockRepo.SetExists(func(key string) bool {
		return true
	})

	s.mockRepo.SetDelete(func(key string) errorUtils.EntityError {
		return expected
	})

	err := services.UserSessionService.DeleteSession(testToken)

	assert.NotNil(s.T(), err)
	assert.EqualValues(s.T(), expected, err)
}

func (s *UserSessionServiceTestSuite) TestSessionExpired_Success() {
	s.mockRepo.SetExists(func(key string) bool {
		return true
	})

	s.mockRepo.SetGet(func(key string) (*domain.UserSession, errorUtils.EntityError) {
		return testSession, nil
	})

	actual, err := services.UserSessionService.IsSessionExpired(testToken, testTimeNow.Add(time.Minute))

	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	assert.True(s.T(), actual)
}

func (s *UserSessionServiceTestSuite) TestSessionExpired_Failure_ErrorGettingSession() {
	expected := errorUtils.NewInternalServerError("error fetching session")
	s.mockRepo.SetExists(func(key string) bool {
		return true
	})

	s.mockRepo.SetGet(func(key string) (*domain.UserSession, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error fetching session")
	})

	actual, err := services.UserSessionService.IsSessionExpired(testToken, testTimeNow.Add(time.Minute))

	assert.NotNil(s.T(), err)
	assert.True(s.T(), actual)
	assert.Equal(s.T(), expected, err)
}

func (s *UserSessionServiceTestSuite) TestSessionExpired_Failure_SessionNotFound() {
	expected := errorUtils.NewNotFoundError(fmt.Sprintf("token with key %s does not exist", testToken))
	s.mockRepo.SetExists(func(key string) bool {
		return false
	})

	expired, err := services.UserSessionService.IsSessionExpired(testToken, testTimeNow.Add(time.Minute))

	assert.NotNil(s.T(), err)
	assert.True(s.T(), expired)
	assert.Equal(s.T(), expected, err)
}
