package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
)

type UserServiceMockInterface interface {
	SetGetUser(func(uint64) (*domain.User, errorUtils.EntityError))
	SetCreateUser(func(*domain.User) (*domain.User, errorUtils.EntityError))
	SetUpdateUser(func(*domain.User) (*domain.User, errorUtils.EntityError))
	SetDelete(func(uint64) errorUtils.EntityError)
	SetGetAll(func() ([]domain.User, errorUtils.EntityError))
}

type UserServiceMock struct {
	getUserService    func(uint64) (*domain.User, errorUtils.EntityError)
	createUserService func(*domain.User) (*domain.User, errorUtils.EntityError)
	updateUserService func(*domain.User) (*domain.User, errorUtils.EntityError)
	deleteUserService func(uint64) errorUtils.EntityError
	getAllUserService func() ([]domain.User, errorUtils.EntityError)
}

func (u *UserServiceMock) GetUser(id uint64) (*domain.User, errorUtils.EntityError) {
	return u.getUserService(id)
}

func (u *UserServiceMock) CreateUser(user *domain.User) (*domain.User, errorUtils.EntityError) {
	return u.createUserService(user)
}

func (u *UserServiceMock) UpdateUser(user *domain.User) (*domain.User, errorUtils.EntityError) {
	return u.updateUserService(user)
}

func (u *UserServiceMock) DeleteUser(id uint64) errorUtils.EntityError {
	return u.deleteUserService(id)
}

func (u *UserServiceMock) GetAllUsers() ([]domain.User, errorUtils.EntityError) {
	return u.getAllUserService()
}

func (u *UserServiceMock) SetGetUser(f func(uint64) (*domain.User, errorUtils.EntityError)) {
	u.getUserService = f
}

func (u *UserServiceMock) SetCreateUser(f func(*domain.User) (*domain.User, errorUtils.EntityError)) {
	u.createUserService = f
}

func (u *UserServiceMock) SetUpdateUser(f func(*domain.User) (*domain.User, errorUtils.EntityError)) {
	u.updateUserService = f
}

func (u *UserServiceMock) SetDelete(f func(uint64) errorUtils.EntityError) {
	u.deleteUserService = f
}

func (u *UserServiceMock) SetGetAll(f func() ([]domain.User, errorUtils.EntityError)) {
	u.getAllUserService = f
}
