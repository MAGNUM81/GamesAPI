package services

import "GamesAPI/src/utils/authUtils"

var (
	AuthenticationService AuthServiceInterface = &AuthService{}
)

type AuthServiceInterface interface {
	ValidatePassword(plainPassword []byte, hashedPassword string) (bool, error)
}

type AuthService struct {}

func (a *AuthService) ValidatePassword(plainPassword []byte, hashedPassword string) (bool, error) {
	return authUtils.CompareStrings(hashedPassword, plainPassword)
}
