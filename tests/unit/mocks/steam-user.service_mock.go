package mocks

import "GamesAPI/src/domain"

type SteamUserMockInterface interface{
	SetGetUserID(func(string) (string, error))
	SetGetUserOwnedGames(func(string) ([]string, error))
	SetGetGameInfo(func(string)(domain.Game, error))
}

type SteamUserMock struct {
	getUserID func(string) (string, error)
	getUserOwnedGames func(string) ([]string, error)
	getGameInfo func(string) (domain.Game, error)
}

func (s *SteamUserMock) GetUserID(personalURL string) (string, error) {
	return s.getUserID(personalURL)
}

func (s *SteamUserMock) GetUserOwnedGames(userID string)([]string, error){
	return s.getUserOwnedGames(userID)
}

func (s *SteamUserMock) GetGameInfo(gameID string)(domain.Game, error){
	return s.getGameInfo(gameID)
}

func (s *SteamUserMock) SetGetUserID(f func(string) (string, error)) {
	s.getUserID = f
}

func (s *SteamUserMock)SetGetUserOwnedGames(f func(string) ([]string, error)){
	s.getUserOwnedGames = f
}

func (s *SteamUserMock)SetGetGameInfo(f func(string)(domain.Game, error)){
	s.getGameInfo = f
}