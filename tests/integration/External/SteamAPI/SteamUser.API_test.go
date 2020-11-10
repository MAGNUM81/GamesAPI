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