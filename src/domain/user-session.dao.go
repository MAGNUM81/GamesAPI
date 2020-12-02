package domain

import "GamesAPI/src/utils/errorUtils"

var (
	UserSessionRepo = NewUserAuthTokenRepository()
)

type UserSessionRepoInterface interface {
	Get(key string) (*UserSession, errorUtils.EntityError)
	Create(key string, token *UserSession) (*UserSession, errorUtils.EntityError)
	Delete(key string) errorUtils.EntityError
	Exists(key string) bool
}

type userSessionRepo struct {
	repo map[string]*UserSession
}

func (u *userSessionRepo) Get(key string) (*UserSession, errorUtils.EntityError) {
	user := u.repo[key]
	var err errorUtils.EntityError = nil
	if user == nil {
		err = errorUtils.NewNotFoundError("Token does not exist in repository")
	}
	return user, err
}

func (u *userSessionRepo) Create(key string, token *UserSession) (*UserSession, errorUtils.EntityError) {
	u.repo[key] = token
	return token, nil
}

func (u *userSessionRepo) Delete(key string) errorUtils.EntityError {
	u.repo[key] = nil
	return nil
}

func (u *userSessionRepo) Exists(key string) bool {
	return u.repo[key] != nil
}

func NewUserAuthTokenRepository() UserSessionRepoInterface {
	return &userSessionRepo{repo: map[string]*UserSession{}}
}
