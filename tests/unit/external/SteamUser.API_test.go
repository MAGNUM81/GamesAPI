package external

import (
	"GamesAPI/src/External/Steam"
	"GamesAPI/src/domain"
	"GamesAPI/tests/unit/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SteamUserAPITestSuite struct {
	suite.Suite
	mock mocks.SteamUserMockInterface
}

func TestSteamUserAPITestSuite(t *testing.T) {
	suite.Run(t, new(SteamUserAPITestSuite))
}

func (s *SteamUserAPITestSuite) SetupSuite() {
	mock := &mocks.SteamUserMock{}
	s.mock = mock
	Steam.ExternalSteamUserService = mock
}

func (s *SteamUserAPITestSuite) TestGetSteamUserID_Success() {
	s.mock.SetGetUserID(func(personalURL string) (string, error) {
		return "76561197960287939", nil
	})
	t := s.T()
	steamUserID, err := Steam.ExternalSteamUserService.GetUserID("gabelogannewell")
	assert.Nil(t, err)
	assert.EqualValues(t, "76561197960287939", steamUserID)
}

func (s *SteamUserAPITestSuite) TestGetSteamUserID_Fail() {
	s.mock.SetGetUserID(func(personalURL string) (string, error) {
		return "", nil
	})
	t := s.T()
	steamUserID, err := Steam.ExternalSteamUserService.GetUserID("invalidUserURL")
	assert.Nil(t, err)
	assert.EqualValues(t, "", steamUserID)
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_Success() {
	s.mock.SetGetUserOwnedGames(func(personalURL string) ([]string, error) {
		return []string{"44", "22"}, nil
	})
	steamUserID := "76561198017133337"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, "44", steamGamesIDs[0])
	assert.EqualValues(t, "22", steamGamesIDs[1])
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_OwnesNoGames() {
	s.mock.SetGetUserOwnedGames(func(personalURL string) ([]string, error) {
		return []string{}, nil
	})
	steamUserID := "76561197960287930"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(steamGamesIDs))
}

func (s *SteamUserAPITestSuite) TestGetSteamUserOwnedGames_BadUserID() {
	s.mock.SetGetUserOwnedGames(func(personalURL string) ([]string, error) {
		return []string{}, nil
	})
	steamUserID := "thishavenochanceofbeingarealsteamid1324567899876544321"
	steamGamesIDs, err := Steam.ExternalSteamUserService.GetUserOwnedGames(steamUserID)
	t := s.T()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(steamGamesIDs))
}

func (s *SteamUserAPITestSuite) TestGetSteamGame_Success(){
	s.mock.SetGetGameInfo(func(gameID string)(domain.Game,error){
		return domain.Game{
			ID:          0,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			DeletedAt:   nil,
			Title:       "NieR:Automata™",
			Developer:   "Square Enix | PlatinumGames Inc.",
			Publisher:   "Square Enix",
			ReleaseDate: time.Date(2017,time.March,17,0,0,0,0,time.UTC),
		}, nil
	})
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
	s.mock.SetGetGameInfo(func(gameID string)(domain.Game,error){
		return domain.Game{
			ID:          0,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			DeletedAt:   nil,
			Title:       "PAYDAY 2",
			Developer:   "OVERKILL - a Starbreeze Studio.",
			Publisher:   "Starbreeze Publishing AB",
			ReleaseDate: time.Date(2013,time.August,13,0,0,0,0,time.UTC),
		}, nil
	})
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
	s.mock.SetGetGameInfo(func(gameID string)(domain.Game,error){
		return domain.Game{}, errors.New("bad Game ID")
	})
	gameID := "65465156435"
	gameInfo, err := Steam.ExternalSteamUserService.GetGameInfo(gameID)
	t := s.T()
	assert.NotNil(t, err)
	assert.EqualValues(t, domain.Game{}, gameInfo)

}