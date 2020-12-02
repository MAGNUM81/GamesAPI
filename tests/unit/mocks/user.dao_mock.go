package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
)

type UserRepoMockInterface interface {
	SetGetUserDomain(func(uint64) (*domain.User, errorUtils.EntityError))
	SetCreateUserDomain(func(user *domain.User) (*domain.User, errorUtils.EntityError))
	SetUpdateUserDomain(func(user *domain.User) (*domain.User, errorUtils.EntityError))
	SetDeleteUserDomain(func(id uint64) errorUtils.EntityError)
	SetGetAllUserDomain(func() ([]domain.User, errorUtils.EntityError))
}

type UserRepoMock struct {
	getUserDomain     func(id uint64) (*domain.User, errorUtils.EntityError)
	createUserDomain  func(user *domain.User) (*domain.User, errorUtils.EntityError)
	updateUserDomain  func(user *domain.User) (*domain.User, errorUtils.EntityError)
	deleteUserDomain  func(id uint64) errorUtils.EntityError
	getAllUsersDomain func() ([]domain.User, errorUtils.EntityError)
}

//UserRepoMockInterface implementation, so we can swap the methods around and get the desired behavior from the repository
func (m *UserRepoMock) SetGetUserDomain(f func(uint64) (*domain.User, errorUtils.EntityError)) {
	m.getUserDomain = f
}

func (m *UserRepoMock) SetCreateUserDomain(f func(user *domain.User) (*domain.User, errorUtils.EntityError)) {
	m.createUserDomain = f
}

func (m *UserRepoMock) SetUpdateUserDomain(f func(user *domain.User) (*domain.User, errorUtils.EntityError)) {
	m.updateUserDomain = f
}

func (m *UserRepoMock) SetDeleteUserDomain(f func(id uint64) errorUtils.EntityError) {
	m.deleteUserDomain = f
}

func (m *UserRepoMock) SetGetAllUserDomain(f func() ([]domain.User, errorUtils.EntityError)) {
	m.getAllUsersDomain = f
}

//UserRepoInterface implementation (redirects all calls to the swappable methods)
func (m *UserRepoMock) Get(id uint64) (*domain.User, errorUtils.EntityError) {
	return m.getUserDomain(id)
}
func (m *UserRepoMock) Create(msg *domain.User) (*domain.User, errorUtils.EntityError) {
	return m.createUserDomain(msg)
}
func (m *UserRepoMock) Update(msg *domain.User) (*domain.User, errorUtils.EntityError) {
	return m.updateUserDomain(msg)
}
func (m *UserRepoMock) Delete(id uint64) errorUtils.EntityError {
	return m.deleteUserDomain(id)
}
func (m *UserRepoMock) GetAll() ([]domain.User, errorUtils.EntityError) {
	return m.getAllUsersDomain()
}
func (m *UserRepoMock) Initialize(_ *gorm.DB) {}
