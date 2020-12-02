package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
)

var (
	UsersService UsersServiceInterface = &usersService{}
)

type usersService struct{}

type UsersServiceInterface interface {
	GetUser(uint64) (*domain.User, errorUtils.EntityError)
	CreateUser(*domain.User) (*domain.User, errorUtils.EntityError)
	UpdateUser(*domain.User) (*domain.User, errorUtils.EntityError)
	DeleteUser(uint64) errorUtils.EntityError
	GetAllUsers() ([]domain.User, errorUtils.EntityError)
}

func (u usersService) GetUser(userId uint64) (*domain.User, errorUtils.EntityError) {
	user, err := domain.UserRepo.Get(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u usersService) CreateUser(user *domain.User) (*domain.User, errorUtils.EntityError) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user, err := domain.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u usersService) UpdateUser(user *domain.User) (*domain.User, errorUtils.EntityError) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	current, err := domain.UserRepo.Get(user.ID)
	if err != nil {
		return nil, err
	}
	current.Email = user.Email
	current.Name = user.Name

	updatedUser, err := domain.UserRepo.Update(current)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (u usersService) DeleteUser(userId uint64) errorUtils.EntityError {
	user, err := domain.UserRepo.Get(userId)
	if err != nil {
		return err
	}

	deleteErr := domain.UserRepo.Delete(user.ID)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (u usersService) GetAllUsers() ([]domain.User, errorUtils.EntityError) {
	users, err := domain.UserRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}
