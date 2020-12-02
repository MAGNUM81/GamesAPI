package mocks

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
)

type UserRoleServiceMockInterface interface {
	SetGetRole(f func(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError))
	SetGetRolesByUserID(f func(userId uint64) ([]domain.UserRole, errorUtils.EntityError))
	SetGetRolesByRoleName(f func(roleName string) ([]domain.UserRole, errorUtils.EntityError))
	SetCreateRole(f func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError))
	SetUpdateRole(f func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError))
	SetDeleteRole(f func(roleId uint64) errorUtils.EntityError)
	SetGetAllRoles(f func() ([]domain.UserRole, errorUtils.EntityError))
}

type UserRoleMock struct {
	getRole            func(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError)
	getRolesByUserID   func(userId uint64) ([]domain.UserRole, errorUtils.EntityError)
	getRolesByRoleName func(roleName string) ([]domain.UserRole, errorUtils.EntityError)
	createRole         func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)
	updateRole         func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)
	deleteRole         func(roleId uint64) errorUtils.EntityError
	getAllRoles        func() ([]domain.UserRole, errorUtils.EntityError)
}

func (u *UserRoleMock) SetGetRole(f func(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError)) {
	u.getRole = f
}

func (u *UserRoleMock) SetGetRolesByUserID(f func(userId uint64) ([]domain.UserRole, errorUtils.EntityError)) {
	u.getRolesByUserID = f
}

func (u *UserRoleMock) SetGetRolesByRoleName(f func(roleName string) ([]domain.UserRole, errorUtils.EntityError)) {
	u.getRolesByRoleName = f
}

func (u *UserRoleMock) SetCreateRole(f func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)) {
	u.createRole = f
}

func (u *UserRoleMock) SetUpdateRole(f func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)) {
	u.updateRole = f
}

func (u *UserRoleMock) SetDeleteRole(f func(roleId uint64) errorUtils.EntityError) {
	u.deleteRole = f
}

func (u *UserRoleMock) SetGetAllRoles(f func() ([]domain.UserRole, errorUtils.EntityError)) {
	u.getAllRoles = f
}

func (u *UserRoleMock) GetRole(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError) {
	return u.getRole(userRoleId)
}

func (u *UserRoleMock) GetRolesByUserID(userId uint64) ([]domain.UserRole, errorUtils.EntityError) {
	return u.getRolesByUserID(userId)
}

func (u *UserRoleMock) GetRolesByRoleName(roleName string) ([]domain.UserRole, errorUtils.EntityError) {
	return u.getRolesByRoleName(roleName)
}

func (u *UserRoleMock) CreateRole(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
	return u.createRole(role)
}

func (u *UserRoleMock) UpdateRole(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
	return u.updateRole(role)
}

func (u *UserRoleMock) DeleteRole(roleId uint64) errorUtils.EntityError {
	return u.deleteRole(roleId)
}

func (u *UserRoleMock) GetAllRoles() ([]domain.UserRole, errorUtils.EntityError) {
	return u.getAllRoles()
}
