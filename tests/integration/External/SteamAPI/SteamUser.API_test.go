package SteamAPI

import (
	"GamesAPI/src/External/Steam"
	"GamesAPI/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SteamUserAPITestSuite struct {
	suite.Suite
}

func TestSteamUserAPITestSuite(t *testing.T){
	suite.Run(t, new(SteamUserAPITestSuite))
}

func (s *SteamUserAPITestSuite) SetupSuite(){
	integration.SimulateEnv()
}
 
func (s *SteamUserAPITestSuite) BeforeTest(_, _ string){
}

func (s *SteamUserAPITestSuite) TestGetSteamUserID_Success(){
	steamUserURL := "gabelogannewell"
	steamUserID, err := Steam.ExternalSteamUserService.GetUserID(steamUserURL)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, "76561197960287930", steamUserID)
}

//Warning, this test is valid until someone create this as a valid UserURL
func (s *SteamUserAPITestSuite) TestGetSteamUserID_BadUserURL(){
	steamUserURL := "gabelogannewell6584968746541654156"
	steamUserID, err := Steam.ExternalSteamUserService.GetUserID(steamUserURL)
	t := s.T()
	assert.EqualValues(t, "", steamUserID)
	assert.NotNil(t, err)
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_Success(){
	steamUserID := "76561198017133337"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, "2100", steamGamesIDs[0])
	assert.EqualValues(t, "2130", steamGamesIDs[1])
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_OwnesNoGames(){
	steamUserID := "76561197960287930"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(steamGamesIDs))
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_BadUserID(){
	steamUserID := "thishavenochanceofbeingarealsteamid1324567899876544321"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(steamGamesIDs))
}