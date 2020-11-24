package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
)

type UserSessionRepoMockInterface interface {
	SetGet(   func(key string) (*domain.UserAuthToken, errorUtils.EntityError))
	SetCreate(func(key string, token *domain.UserAuthToken)(*domain.UserAuthToken, errorUtils.EntityError))
	SetDelete(func(key string) errorUtils.EntityError)
	SetExists(func(key string) bool)
}

type UserSessionRepoMock struct {
	get func(key string) (*domain.UserAuthToken, errorUtils.EntityError)
	create func(key string, token *domain.UserAuthToken)(*domain.UserAuthToken, errorUtils.EntityError)
	delete func(key string) errorUtils.EntityError
	exists func(key string) bool
}

func (m UserSessionRepoMock) SetGet(f func(key string) (*domain.UserAuthToken, errorUtils.EntityError))  {
	m.get = f
}

func (m UserSessionRepoMock) SetGet(f func(key string) (*domain.UserAuthToken, errorUtils.EntityError))  {
	m.create = f
}

func (m UserSessionRepoMock) SetGet(f func(key string) (*domain.UserAuthToken, errorUtils.EntityError))  {
	m.delete = f
}

func (m UserSessionRepoMock) SetGet(f func(key string) (*domain.UserAuthToken, errorUtils.EntityError))  {
	m.exists = f
}
