package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
)

var (
	UserRoleService UserRoleServiceInterface = &userRoleService{}
)

type userRoleService struct{}

type UserRoleServiceInterface interface {
	GetRole(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError)
	GetRolesByUserID(userId uint64) ([]domain.UserRole, errorUtils.EntityError)
	GetRolesByRoleName(roleName string) ([]domain.UserRole, errorUtils.EntityError)
	CreateRole(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)
	UpdateRole(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError)
	DeleteRole(roleId uint64) errorUtils.EntityError
	GetAllRoles() ([]domain.UserRole, errorUtils.EntityError)
}

func (u userRoleService) GetRole(userRoleId uint64) (*domain.UserRole, errorUtils.EntityError) {
	role, err := domain.UserRoleRepo.GetByID(userRoleId)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (u userRoleService) GetRolesByUserID(userId uint64) ([]domain.UserRole, errorUtils.EntityError) {
	roles, err := domain.UserRoleRepo.GetByUserID(userId)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u userRoleService) GetRolesByRoleName(roleName string) ([]domain.UserRole, errorUtils.EntityError) {
	roles, err := domain.UserRoleRepo.GetByRole(roleName)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u userRoleService) CreateRole(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
	role, err := domain.UserRoleRepo.Create(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (u userRoleService) UpdateRole(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
	current, err := domain.UserRoleRepo.GetByID(role.ID)
	if err != nil {
		return nil, err
	}
	current.UserID = role.UserID
	current.Name = role.Name

	updatedRole, err := domain.UserRoleRepo.Update(current)
	if err != nil {
		return nil, err
	}
	return updatedRole, nil
}

func (u userRoleService) DeleteRole(roleId uint64) errorUtils.EntityError {
	current, err := domain.UserRoleRepo.GetByID(roleId)
	if err != nil {
		return err
	}
	deleteErr := domain.UserRoleRepo.Delete(current.ID)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (u userRoleService) GetAllRoles() ([]domain.UserRole, errorUtils.EntityError) {
	roles, err := domain.UserRoleRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
