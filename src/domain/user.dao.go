package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
)

type UserRepoInterface interface {
	Get(uint642 uint64) (*User, errorUtils.EntityError)
	Create(*User) (*User, errorUtils.EntityError)
	Update(*User) (*User, errorUtils.EntityError)
	Delete(uint64) errorUtils.EntityError
	GetAll() ([]User, errorUtils.EntityError)
	Initialize(*gorm.DB)
}

var (
	UserRepo UserRepoInterface = &userRepo{}
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &userRepo{db: db}
}

func (u *userRepo) Initialize(db *gorm.DB) {
	u.db = db
	db.AutoMigrate(&User{})
}

func (u *userRepo) Get(userId uint64) (*User, errorUtils.EntityError) {
	var user User
	if err := u.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	return &user, nil
}

func (u *userRepo) Create(user *User) (*User, errorUtils.EntityError) {
	if dbc := u.db.Create(user); dbc.Error != nil {
		return nil, errorUtils.NewInternalServerError(dbc.Error.Error())
	}
	return user, nil
}

func (u *userRepo) Update(user *User) (*User, errorUtils.EntityError) {
	var found User
	if err := u.db.Where("id = ?", user.ID).First(&found).Error; err != nil {
		return nil, errorUtils.NewNotFoundError(err.Error())
	}
	u.db.Save(*user)
	return user, nil
}

func (u *userRepo) Delete(userId uint64) errorUtils.EntityError {
	var user User
	if err := u.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return errorUtils.NewNotFoundError(err.Error())
	}
	dbc := u.db.Delete(&user)
	return errorUtils.NewEntityError(dbc.Error)
}

func (u *userRepo) GetAll() ([]User, errorUtils.EntityError) {
	var users []User
	u.db.Find(&users)
	return users, nil
}