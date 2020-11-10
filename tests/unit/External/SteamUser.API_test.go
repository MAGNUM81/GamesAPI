package External

import (
	"GamesAPI/src/External/Steam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SteamUserMockInterface interface{
	SetGetUserID(func(string) (string, error))
}

type steamUserMock struct {
	getUserID func(string) (string, error)
}

func (s *steamUserMock) GetUserID(personalURL string) (string, error) {
	return s.getUserID(personalURL)
}

func (s *steamUserMock) SetGetUserID(f func(string) (string, error)) {
	s.getUserID = f
}

type SteamUserAPITestSuite struct {
	suite.Suite
	mock SteamUserMockInterface
}

func TestSteamUserAPITestSuite(t *testing.T){
	suite.Run(t, new(SteamUserAPITestSuite))
}

func (s *SteamUserAPITestSuite) SetupSuite() {
	mock := &steamUserMock{}
	s.mock = mock
	Steam.ExternalSteamUserService = mock
}

func (s *SteamUserAPITestSuite) TestGetSteamUserID_Success() {
	s.mock.SetGetUserID(func(personalURL string) (string, error){
		return "76561197960287939", nil
	})
	t := s.T()
	steamUserID, err := Steam.ExternalSteamUserService.GetUserID("gabelogannewell")
	assert.Nil(t, err)
	assert.EqualValues(t, "76561197960287939", steamUserID)
}

func (s *SteamUserAPITestSuite) TestGetSteamUserID_Fail() {
	s.mock.SetGetUserID(func(personalURL string) (string, error){
		return "", nil
	})
	t := s.T()
	steamUserID, err := Steam.ExternalSteamUserService.GetUserID("invalidUserURL")
	assert.Nil(t, err)
	assert.EqualValues(t, "", steamUserID)
}