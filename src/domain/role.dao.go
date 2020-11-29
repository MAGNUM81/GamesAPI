package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
)

var (
	UserRoleRepo UserRoleRepoInterface = &userRoleRepo{}
)

type UserRoleRepoInterface interface {
	GetByID(uint64) (*UserRole, errorUtils.EntityError)
	GetByUserID(uint64) ([]UserRole, errorUtils.EntityError)
	GetByRole(string) ([]UserRole, errorUtils.EntityError)
	Create(*UserRole) (*UserRole, errorUtils.EntityError)
	Update(*UserRole) (*UserRole, errorUtils.EntityError)
	Delete(uint64) errorUtils.EntityError
	GetAll() ([]UserRole, errorUtils.EntityError)
	Initialize(db *gorm.DB)
}

type userRoleRepo struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepoInterface {
	return &userRoleRepo{db: db}
}

func (u *userRoleRepo) Create(role *UserRole) (*UserRole, errorUtils.EntityError) {
	if dbc := u.db.Create(role); dbc.Error != nil {
		return nil, errorUtils.NewInternalServerError(dbc.Error.Error())
	}
	return role, nil
}

func (u *userRoleRepo) Update(role *UserRole) (*UserRole, errorUtils.EntityError) {
	if err := u.db.Where("id = ?", role.ID).First(&role).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	u.db.Save(*role)
	return role, nil
}

func (u *userRoleRepo) Initialize(db *gorm.DB) {
	u.db = db
	db.AutoMigrate(&UserRole{})
}

func (u *userRoleRepo) GetByID(roleId uint64) (*UserRole, errorUtils.EntityError) {
	var userRole UserRole
	if err := u.db.Where("id = ?", roleId).First(&userRole).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	return &userRole, nil
}

func (u *userRoleRepo) GetByUserID(userId uint64) ([]UserRole, errorUtils.EntityError) {
	var userRoles []UserRole
	if err := u.db.Find(&userRoles, "user_id = ?", userId).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	return userRoles, nil
}

func (u *userRoleRepo) GetByRole(roleName string) ([]UserRole, errorUtils.EntityError) {
	var userRoles []UserRole
	if err := u.db.Find(&userRoles, "roleName = ?", roleName).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	return userRoles, nil
}

func (u *userRoleRepo) Delete(roleId uint64) errorUtils.EntityError {
	var role UserRole
	if err := u.db.Where("id = ?", roleId).First(&role).Error; err != nil {
		return errorUtils.NewNotFoundError(err.Error())
	}
	dbc := u.db.Delete(&role)
	return errorUtils.NewEntityError(dbc.Error)
}

func (u *userRoleRepo) GetAll() ([]UserRole, errorUtils.EntityError) {
	var roles []UserRole
	u.db.Find(&roles)
	return roles, nil
}
