package domain

import "GamesAPI/src/utils/errorUtils"

var (
	UserAuthTokenRepo UserAuthTokenRepoInterface = NewUserAuthTokenRepository())

type UserAuthTokenRepoInterface interface {
	Get(key string) (*UserAuthToken, errorUtils.EntityError)
	Create(key string, token *UserAuthToken) (*UserAuthToken, errorUtils.EntityError)
	Delete(key string) errorUtils.EntityError
	Exists(key string) bool
}

type userAuthTokenRepo struct {
	repo map[string]*UserAuthToken
}

func (u userAuthTokenRepo) Get(key string) (*UserAuthToken, errorUtils.EntityError) {
	user := u.repo[key]
	var err errorUtils.EntityError = nil
	if user == nil {
		err = errorUtils.NewNotFoundError("Token does not exist in repository")
	}
	return user, err
}

func (u userAuthTokenRepo) Create(key string, token *UserAuthToken) (*UserAuthToken, errorUtils.EntityError) {
	u.repo[key] = token
	return token, nil
}

func (u userAuthTokenRepo) Delete(key string) errorUtils.EntityError {
	u.repo[key] = nil
	return nil
}

func (u userAuthTokenRepo) Exists(key string) bool {
	return u.repo[key] != nil
}

func NewUserAuthTokenRepository() UserAuthTokenRepoInterface{
	return &userAuthTokenRepo{repo: map[string]*UserAuthToken{}}
}


