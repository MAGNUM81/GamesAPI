package services

import (
	"GamesAPI/src/services"
	"GamesAPI/src/utils/authUtils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthenticationServiceTestSuite struct {
	suite.Suite
}

func TestAuthenticationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationServiceTestSuite))
}

func (s *AuthenticationServiceTestSuite) TestValidatePassword_Success() {
	strPlainPassword := "thisisanicep4ssw0rd"
	hashed, err := authUtils.HashAndSalt([]byte(strPlainPassword))

	assert.NotNil(s.T(), hashed)
	assert.Nil(s.T(), err)

	validates, err := services.AuthenticationService.ValidatePassword([]byte(strPlainPassword), hashed)

	assert.True(s.T(), validates)
	assert.Nil(s.T(), err)

}

func (s *AuthenticationServiceTestSuite) TestValidatePassword_Failure() {
	strGoodPlainPassword := "thisisanicep4ssw0rd"
	hashed, err := authUtils.HashAndSalt([]byte(strGoodPlainPassword))

	assert.NotNil(s.T(), hashed)
	assert.Nil(s.T(), err)

	strBadPlainPassword := "p4sswordanicethisis"

	validates, err := services.AuthenticationService.ValidatePassword([]byte(strBadPlainPassword), hashed)

	assert.False(s.T(), validates)
	assert.NotNil(s.T(), err)
}
