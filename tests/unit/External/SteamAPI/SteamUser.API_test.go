package SteamAPI

import (
	"GamesAPI/src/External/Steam"
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
}

func (s *SteamUserAPITestSuite) BeforeTest(_, _ string){
}

func (s *SteamUserAPITestSuite) TestGetSteamUserID_Success(){
	steamUserURL := "gabelogannewell"
	steamUserID := Steam.GetUserID(steamUserURL)
	t := s.T()
	assert.EqualValues(t, "76561197960287930", steamUserID)
}

//Warning, this test is valid until someone create this as a valid UserURL
func (s *SteamUserAPITestSuite) TestGetSteamUserID_BadUserURL(){
	steamUserURL := "gabelogannewell6584968746541654156"
	steamUserID := Steam.GetUserID(steamUserURL)
	t := s.T()
	assert.EqualValues(t, "No match", steamUserID)
}