package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"time"
)

type UserSessionServiceMockInterface interface {
	SetGenerateSessionToken(f func(userId uint64, expireAt time.Time) (string, error))
	SetCreateSession(f func(token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError))
	SetIsSessionExpired(f func(key string, currentTime time.Time) (bool, errorUtils.EntityError))
	SetExistsSession(f func(key string) bool)
	SetDeleteSession(f func(key string) errorUtils.EntityError)
}

type UserSessionServiceMock struct {
	generateSessionToken func(userId uint64, expireAt time.Time) (string, error)
	createSession func(token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError)
	isSessionExpired func(key string, currentTime time.Time) (bool, errorUtils.EntityError)
	existsSession func(key string) bool
	deleteSession func(key string) errorUtils.EntityError
}

func (m *UserSessionServiceMock) SetGenerateSessionToken(f func(userId uint64, expireAt time.Time) (string, error)) {
	m.generateSessionToken = f
}
