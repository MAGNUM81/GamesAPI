package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/tests/unit/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
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

func (s *UserSessionServiceTestSuite) TestCreateToken_Success() {

}

func (s *UserSessionServiceTestSuite) TestCreateToken_Failure(){

}

func (s *UserSessionServiceTestSuite) TestGetToken_Success(){

}

func (s *UserSessionServiceTestSuite) TestGetToken_Failure(){

}

func (s *UserSessionServiceTestSuite) TestDeleteToken_Success(){

}

func (s *UserSessionServiceTestSuite) TestDeleteToken_Failure(){

}



