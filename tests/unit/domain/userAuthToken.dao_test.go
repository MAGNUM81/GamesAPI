package domain

import (
	"GamesAPI/src/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

var (
	now = time.Now().Unix()
)

type UATS struct {
	suite.Suite
	userAuthToken *domain.UserAuthToken
}

func TestUserAuthTokenTestSuite(t *testing.T) {
	suite.Run(t, new(UATS))
}

func (s *UATS) BeforeTest(_, _ string) {
	domain.UserAuthTokenRepo = domain.NewUserAuthTokenRepository()
	s.userAuthToken = &domain.UserAuthToken{
		Token:     "abcdef",
		UserId:    1,
		ExpiresAt: now,
	}
}

func (s *UATS) TestRepo_Get_Empty() {
	key := "abs"
	u, err := domain.UserAuthTokenRepo.Get(key)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), u)
	assert.Equal(s.T(), http.StatusNotFound, err.Status())
	assert.Equal(s.T(), "Token does not exist in repository", err.Message())
}

func (s *UATS) TestRepo_Get_Exists() {
	key := "abcdef"
	_, _ = domain.UserAuthTokenRepo.Create(key, s.userAuthToken)
	u, err := domain.UserAuthTokenRepo.Get(key)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), u)
	assert.Equal(s.T(), s.userAuthToken, u)
}

func (s *UATS) TestRepo_Get_NotExists() {
	key := "abcdef"
	otherKey := "tyuio"
	_, _ = domain.UserAuthTokenRepo.Create(otherKey, s.userAuthToken)
	u, err := domain.UserAuthTokenRepo.Get(key)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), u)
	assert.Equal(s.T(), "Token does not exist in repository", err.Message())
}

func (s *UATS) TestRepo_Create_New() {
	key := "abcdef"
	authToken := s.userAuthToken
	u, err := domain.UserAuthTokenRepo.Create(key, authToken)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), u)
	assert.Equal(s.T(), s.userAuthToken, u)

}

func (s *UATS) TestRepo_Create_KeyExists() {
	key := "abcdef"
	authToken := s.userAuthToken
	_, _ = domain.UserAuthTokenRepo.Create(key, authToken)
	u, err := domain.UserAuthTokenRepo.Create(key, authToken)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), u)
	assert.Equal(s.T(), s.userAuthToken, u)
}

func (s *UATS) TestRepo_Delete_Exists() {
	key := "bji"
	authToken := s.userAuthToken
	_, _ = domain.UserAuthTokenRepo.Create(key, authToken)

	err := domain.UserAuthTokenRepo.Delete(key)
	assert.Nil(s.T(), err)

	u, _ := domain.UserAuthTokenRepo.Get(key)
	assert.Nil(s.T(), u)

}

func (s *UATS) TestRepo_Delete_NotExists() {
	key := "bji"
	err := domain.UserAuthTokenRepo.Delete(key)
	assert.Nil(s.T(), err)
	u, _ := domain.UserAuthTokenRepo.Get(key)
	assert.Nil(s.T(), u)
}

func (s *UATS) TestRepo_Exists_Exists() {
	key := "bji"
	authToken := s.userAuthToken
	_, _ = domain.UserAuthTokenRepo.Create(key, authToken)
	assert.True(s.T(), domain.UserAuthTokenRepo.Exists(key))
}

func (s *UATS) TestRepo_Exists_NotExists() {
	key := "bji"
	assert.False(s.T(), domain.UserAuthTokenRepo.Exists(key))
}