package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"time"
)

type UserSessionServiceMockInterface interface {
	SetGetSession(f func(key string) (*domain.UserSession, errorUtils.EntityError))
	SetGenerateSessionToken(f func(userId uint64, expireAt time.Time) (string, error))
	SetCreateSession(f func(token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError))
	SetIsSessionExpired(f func(key string, currentTime time.Time) (bool, errorUtils.EntityError))
	SetExistsSession(f func(key string) bool)
	SetDeleteSession(f func(key string) errorUtils.EntityError)
}

type UserSessionServiceMock struct {
	getSession func(key string) (*domain.UserSession, errorUtils.EntityError)
	generateSessionToken func(userId uint64, expireAt time.Time) (string, error)
	createSession        func(token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError)
	isSessionExpired     func(key string, currentTime time.Time) (bool, errorUtils.EntityError)
	existsSession        func(key string) bool
	deleteSession        func(key string) errorUtils.EntityError
}

func (m *UserSessionServiceMock) GetSession(key string) (*domain.UserSession, errorUtils.EntityError) {
	return m.getSession(key)
}

func (m *UserSessionServiceMock) CreateSession(token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError) {
	return m.createSession(token)
}

func (m *UserSessionServiceMock) ExistsSession(token string) bool {
	return m.existsSession(token)
}

func (m *UserSessionServiceMock) DeleteSession(token string) errorUtils.EntityError {
	return m.deleteSession(token)
}

func (m *UserSessionServiceMock) IsSessionExpired(key string, currentTime time.Time) (bool, errorUtils.EntityError) {
	return m.isSessionExpired(key, currentTime)
}

func (m *UserSessionServiceMock) GenerateSessionToken(userId uint64, expireAt time.Time) (string, error) {
	return m.generateSessionToken(userId, expireAt)
}

func (m *UserSessionServiceMock) SetGetSession(f func(key string) (*domain.UserSession, errorUtils.EntityError)) {
	m.getSession = f
}

func (m *UserSessionServiceMock) SetCreateSession(f func(token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError)) {
	m.createSession = f
}

func (m *UserSessionServiceMock) SetIsSessionExpired(f func(key string, currentTime time.Time) (bool, errorUtils.EntityError)) {
	m.isSessionExpired = f
}

func (m *UserSessionServiceMock) SetExistsSession(f func(key string) bool) {
	m.existsSession = f
}

func (m *UserSessionServiceMock) SetDeleteSession(f func(key string) errorUtils.EntityError) {
	m.deleteSession = f
}

func (m *UserSessionServiceMock) SetGenerateSessionToken(f func(userId uint64, expireAt time.Time) (string, error)) {
	m.generateSessionToken = f
}
