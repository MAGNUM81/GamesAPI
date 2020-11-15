package External

import (
	"GamesAPI/src/External/Steam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SteamUserMockInterface interface{
	SetGetUserID(func(string) (string, error))
	SetGetUserOwnedGames(func(string) ([]string, error))
}

type steamUserMock struct {
	getUserID func(string) (string, error)
	getUserOwnedGames func(string) ([]string, error)
}

func (s *steamUserMock) GetUserID(personalURL string) (string, error) {
	return s.getUserID(personalURL)
}

func (s *steamUserMock) GetUserOwnedGames(userID string)([]string, error){
	return s.getUserOwnedGames(userID)
}

func (s *steamUserMock) SetGetUserID(f func(string) (string, error)) {
	s.getUserID = f
}

func (s *steamUserMock)SetGetUserOwnedGames(f func(string) ([]string, error)){
	s.getUserOwnedGames = f
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

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_Success(){
	s.mock.SetGetUserOwnedGames(func(personalURL string) ([]string, error){
		return []string{"44","22"}, nil
	})
	steamUserID := "76561198017133337"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, "44", steamGamesIDs[0])
	assert.EqualValues(t, "22", steamGamesIDs[1])
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_OwnesNoGames(){
	s.mock.SetGetUserOwnedGames(func(personalURL string) ([]string, error){
		return []string{}, nil
	})
	steamUserID := "76561197960287930"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(steamGamesIDs))
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_BadUserID(){
	s.mock.SetGetUserOwnedGames(func(personalURL string) ([]string, error){
		return []string{}, nil
	})
	steamUserID := "thishavenochanceofbeingarealsteamid1324567899876544321"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(steamGamesIDs))
}