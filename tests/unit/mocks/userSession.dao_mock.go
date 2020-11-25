package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
)

type UserSessionRepoMockInterface interface {
	SetGet(   func(key string) (*domain.UserSession, errorUtils.EntityError))
	SetCreate(func(key string, token *domain.UserSession)(*domain.UserSession, errorUtils.EntityError))
	SetDelete(func(key string) errorUtils.EntityError)
	SetExists(func(key string) bool)
}

type UserSessionRepoMock struct {
	get func(key string) (*domain.UserSession, errorUtils.EntityError)
	create func(key string, token *domain.UserSession)(*domain.UserSession, errorUtils.EntityError)
	delete func(key string) errorUtils.EntityError
	exists func(key string) bool
}

func (m *UserSessionRepoMock) Get(key string) (*domain.UserSession, errorUtils.EntityError) {
	return m.get(key)
}

func (m *UserSessionRepoMock) Create(key string, token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError) {
	return m.create(key, token)
}

func (m *UserSessionRepoMock) Delete(key string) errorUtils.EntityError {
	return m.delete(key)
}

func (m *UserSessionRepoMock) Exists(key string) bool {
	return m.exists(key)
}

func (m *UserSessionRepoMock) SetCreate(f func(key string, token *domain.UserSession) (*domain.UserSession, errorUtils.EntityError)) {
	m.create = f
}

func (m *UserSessionRepoMock) SetDelete(f func(key string) errorUtils.EntityError) {
	m.delete = f
}

func (m *UserSessionRepoMock) SetExists(f func(key string) bool) {
	m.exists = f
}

func (m *UserSessionRepoMock) SetGet(f func(key string) (*domain.UserSession, errorUtils.EntityError))  {
	m.get = f
}
