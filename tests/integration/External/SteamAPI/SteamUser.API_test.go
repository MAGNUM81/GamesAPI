package SteamAPI

import (
	"GamesAPI/src/External/Steam"
	"GamesAPI/src/domain"
	"GamesAPI/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SteamUserAPITestSuite struct {
	suite.Suite
}

func TestSteamUserAPITestSuite(t *testing.T) {
	suite.Run(t, new(SteamUserAPITestSuite))
}

func (s *SteamUserAPITestSuite) SetupSuite() {
	integration.SimulateEnv()
}

func (s *SteamUserAPITestSuite) BeforeTest(_, _ string) {
}

func (s *SteamUserAPITestSuite) TestGetSteamUserID_Success() {
	steamUserURL := "gabelogannewell"
	steamUserID, err := Steam.ExternalSteamUserService.GetUserID(steamUserURL)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, "76561197960287930", steamUserID)
}

//Warning, this test is valid until someone create this as a valid UserURL
func (s *SteamUserAPITestSuite) TestGetSteamUserID_BadUserURL() {
	steamUserURL := "gabelogannewell6584968746541654156"
	steamUserID, err := Steam.ExternalSteamUserService.GetUserID(steamUserURL)
	t := s.T()
	assert.EqualValues(t, "", steamUserID)
	assert.NotNil(t, err)
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_Success() {
	steamUserID := "76561198017133337"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, "2100", steamGamesIDs[0])
	assert.EqualValues(t, "2130", steamGamesIDs[1])
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_OwnesNoGames() {
	steamUserID := "76561197960287930"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(steamGamesIDs))
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_BadUserID() {
	steamUserID := "thishavenochanceofbeingarealsteamid1324567899876544321"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(steamGamesIDs))
}

func (s *SteamUserAPITestSuite) TestGetSteamGame_Success(){
	gameID := "524220"
	gameInfo, err := Steam.ExternalSteamUserService.GetGameInfo(gameID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, "NieR:Automata™", gameInfo.Title)
	assert.EqualValues(t, "Square Enix | PlatinumGames Inc.", gameInfo.Developer)
	assert.EqualValues(t, "Square Enix", gameInfo.Publisher)
	assert.EqualValues(t, time.Date(2017,time.March,17,0,0,0,0,time.UTC), gameInfo.ReleaseDate)
}

func (s *SteamUserAPITestSuite) TestGetSteamGame_SuccessSecondGame(){
	gameID := "218620"
	gameInfo, err := Steam.ExternalSteamUserService.GetGameInfo(gameID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, "PAYDAY 2", gameInfo.Title)
	assert.EqualValues(t, "OVERKILL - a Starbreeze Studio.", gameInfo.Developer)
	assert.EqualValues(t, "Starbreeze Publishing AB", gameInfo.Publisher)
	assert.EqualValues(t, time.Date(2013,time.August,13,0,0,0,0,time.UTC), gameInfo.ReleaseDate)
}

func (s *SteamUserAPITestSuite) TestGetSteamGame_BadGameID(){
	gameID := "65465156435"
	gameInfo, err := Steam.ExternalSteamUserService.GetGameInfo(gameID)
	t := s.T()
	assert.NotNil(t, err)
	assert.EqualValues(t, domain.Game{}, gameInfo)

}

