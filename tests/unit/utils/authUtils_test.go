package utils

import (
	"GamesAPI/src/utils/authUtils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type AuthUtilsTestSuite struct {
	suite.Suite
}

func TestAuthUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUtilsTestSuite))
}

func (s *AuthUtilsTestSuite) TestJwtSymmetric() {
	var userId uint64 = 1
	expiresAt := time.Now().Unix()
	jwtToken, errCreate := authUtils.JwtCreate(userId, expiresAt)
	assert.Nil(s.T(), errCreate)

	deserializedJwtToken, errDecode := authUtils.JwtDecode(jwtToken)
	claims := deserializedJwtToken.Claims.(*authUtils.UserClaims)
	assert.Nil(s.T(), errDecode)
	assert.Equal(s.T(), userId, claims.UserId)
	assert.Equal(s.T(), expiresAt, claims.ExpiresAt)
}

func (s *AuthUtilsTestSuite) TestPasswordHashSymmetric() {
	password := "this is the best password"
	pwdBytes := []byte(password)
	hash, _ := authUtils.HashAndSalt(pwdBytes)

	areEqual, _ := authUtils.CompareStrings(hash, pwdBytes)
	require.True(s.T(), areEqual)
}


