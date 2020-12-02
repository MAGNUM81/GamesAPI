package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
)

//1. create RepoMock interface
//	put Set(func) methods in there
type UserRoleRepoMockInterface interface {
	SetGetRole(func(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError))
	SetGetRolesByUserID(func(userId uint64) ([]domain.UserRole, errorUtils.EntityError))
	SetGetRolesByRoleName(func(roleName string) ([]domain.UserRole, errorUtils.EntityError))
	SetCreateRole(func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError))
	SetUpdateRole(func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError))
	SetDeleteRole(func(roleId uint64) errorUtils.EntityError)
	SetGetAllRoles(func() ([]domain.UserRole, errorUtils.EntityError))
}

//2. create RepoMock struct with Repo Interface methods as members

type UserRoleRepoMock struct {
	getRole            func(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError)
	getRolesByUserID   func(userId uint64) ([]domain.UserRole, errorUtils.EntityError)
	getRolesByRoleName func(roleName string) ([]domain.UserRole, errorUtils.EntityError)
	createRole         func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)
	updateRole         func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)
	deleteRole         func(roleId uint64) errorUtils.EntityError
	getAllRoles        func() ([]domain.UserRole, errorUtils.EntityError)
}

//3. Implement RepoMock interface in RepoMock struct
// 	assign correct member to passed function
//  these methods can be auto generated once struct constructor is called

func (u *UserRoleRepoMock) SetGetRole(f func(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError)) {
	u.getRole = f
}

func (u *UserRoleRepoMock) SetGetRolesByUserID(f func(userId uint64) ([]domain.UserRole, errorUtils.EntityError)) {
	u.getRolesByUserID = f
}

func (u *UserRoleRepoMock) SetGetRolesByRoleName(f func(roleName string) ([]domain.UserRole, errorUtils.EntityError)) {
	u.getRolesByRoleName = f
}

func (u *UserRoleRepoMock) SetCreateRole(f func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)) {
	u.createRole = f
}

func (u *UserRoleRepoMock) SetUpdateRole(f func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)) {
	u.updateRole = f
}

func (u *UserRoleRepoMock) SetDeleteRole(f func(roleId uint64) errorUtils.EntityError) {
	u.deleteRole = f
}

func (u *UserRoleRepoMock) SetGetAllRoles(f func() ([]domain.UserRole, errorUtils.EntityError)) {
	u.getAllRoles = f
}

//4. Implement Repo Interface in RepoMock struct
//	redirect all calls to mock method
//  these methods can be auto-generated once struct constructor is called

func (u *UserRoleRepoMock) GetByID(u2 uint64) (*domain.UserRole, errorUtils.EntityError) {
	return u.getRole(u2)
}

func (u *UserRoleRepoMock) GetByUserID(u2 uint64) ([]domain.UserRole, errorUtils.EntityError) {
	return u.getRolesByUserID(u2)
}

func (u *UserRoleRepoMock) GetByRole(s string) ([]domain.UserRole, errorUtils.EntityError) {
	return u.getRolesByRoleName(s)
}

func (u *UserRoleRepoMock) Create(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
	return u.createRole(role)
}

func (u *UserRoleRepoMock) Update(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
	return u.updateRole(role)
}

func (u *UserRoleRepoMock) Delete(u2 uint64) errorUtils.EntityError {
	return u.deleteRole(u2)
}

func (u *UserRoleRepoMock) GetAll() ([]domain.UserRole, errorUtils.EntityError) {
	return u.getAllRoles()
}

func (u *UserRoleRepoMock) Initialize(_ *gorm.DB) {}
