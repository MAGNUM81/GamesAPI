package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"
)

var (
	UserSessionService UserSessionServiceInterface = &userSessionService{}
)

type UserSessionServiceInterface interface {
	CreateSession(token *domain.UserAuthToken) (*domain.UserAuthToken, errorUtils.EntityError)
	ExistsSession(token string) bool
	DeleteSession(token string) errorUtils.EntityError
	GenerateSessionToken(userId uint64, expireAt time.Time) (string, error)
}

type userSessionService struct {}

func (u userSessionService) GenerateSessionToken(userId uint64, expireAt time.Time) (string, error) {
	h := fnv.New32a()
	// add both values as bytes to a buffer big enough to contain them
	buf := make([]byte, binary.MaxVarintLen64 + binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, userId)
	m := binary.PutVarint(buf, expireAt.UnixNano())
	//take the resulting slice
	b := buf[:(n+m)]
	//hash it
	_, err := h.Write(b)
	return strconv.Itoa(int(h.Sum32())), err
}

func (u userSessionService) CreateSession(token *domain.UserAuthToken) (*domain.UserAuthToken, errorUtils.EntityError) {
	if err := token.Validate(); err != nil {
		return nil, err
	}

	if domain.UserAuthTokenRepo.Exists(token.Token){
		return nil, errorUtils.NewUnprocessableEntityError(fmt.Sprintf("token with key %s already exists", token.Token))
	}

	ret, err := domain.UserAuthTokenRepo.Create(token.Token, token)
	if err != nil {
		return nil, err
	}

	return ret, err
}

func (u userSessionService) ExistsSession(key string) bool {
	return domain.UserAuthTokenRepo.Exists(key)
}

func (u userSessionService) DeleteSession(key string) errorUtils.EntityError {
	if !domain.UserAuthTokenRepo.Exists(key) {
		return errorUtils.NewNotFoundError(fmt.Sprintf("token with key %s does not exist", key))
	}
	return domain.UserAuthTokenRepo.Delete(key)
}
