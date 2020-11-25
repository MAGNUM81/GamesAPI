package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/tests/unit/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type UserSessionServiceTestSuite struct {
	suite.Suite
	mockRepo mocks.UserSessionRepoMockInterface
}

func TestUserSessionServiceTestSuite(t  *testing.T) {
	suite.Run(t, new(UserSessionServiceTestSuite))
}

func (s *UserSessionServiceTestSuite) SetupSuite() {
	mock := &mocks.UserSessionRepoMock{}
	s.mockRepo = mock
	domain.UserSessionRepo = mock
}

var (
	testTime = time.Now()
	testUserId = uint64(42)
)

func (s *UserSessionServiceTestSuite) TestGenerateSessionToken_Success() {
	
}

func (s *UserSessionServiceTestSuite) TestGenerateSessionToken_Failure() {

}

func (s *UserSessionServiceTestSuite) TestCreateSession_Success() {
	s.mockRepo.SetExists(func(key string) bool {
		return true
	})
}

func (s *UserSessionServiceTestSuite) TestCreateSession_Failure(){
	s.mockRepo.SetExists(func(key string) bool {
		return false
	})
}

func (s *UserSessionServiceTestSuite) TestGetSession_Success(){

}

func (s *UserSessionServiceTestSuite) TestGetSession_Failure(){

}

func (s *UserSessionServiceTestSuite) TestDeleteSession_Success(){

}

func (s *UserSessionServiceTestSuite) TestDeleteSession_Failure(){

}



