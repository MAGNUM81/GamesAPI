package controllers

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/router"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UserServiceMockInterface interface {
	SetGetUser(func(uint64) (*domain.User, errorUtils.EntityError))
	SetCreateUser(func(*domain.User) (*domain.User, errorUtils.EntityError))
	SetUpdateUser(func(*domain.User) (*domain.User, errorUtils.EntityError))
	SetDelete(func(uint64) errorUtils.EntityError)
	SetGetAll(func() ([]domain.User, errorUtils.EntityError))
}

type userServiceMock struct {
	services.UsersServiceInterface
	getUserService func(uint64) (*domain.User, errorUtils.EntityError)
	createUserService func(*domain.User) (*domain.User, errorUtils.EntityError)
	updateUserService func(*domain.User) (*domain.User, errorUtils.EntityError)
	deleteUserService func(uint64) errorUtils.EntityError
	getAllUserService func() ([]domain.User, errorUtils.EntityError)
}

func (u *userServiceMock) GetUser(id uint64) (*domain.User, errorUtils.EntityError) {
	return u.getUserService(id)
}

func (u *userServiceMock) CreateUser(user *domain.User) (*domain.User, errorUtils.EntityError) {
	return u.createUserService(user)
}

func (u *userServiceMock) UpdateUser(user *domain.User) (*domain.User, errorUtils.EntityError) {
	return u.updateUserService(user)
}

func (u *userServiceMock) DeleteUser(id uint64) errorUtils.EntityError {
	return u.deleteUserService(id)
}

func (u *userServiceMock) GetAllUsers() ([]domain.User, errorUtils.EntityError) {
	return u.getAllUserService()
}

func (u *userServiceMock) SetGetUser(f func(uint64) (*domain.User, errorUtils.EntityError)) {
	u.getUserService = f
}

func (u *userServiceMock) SetCreateUser(f func(*domain.User) (*domain.User, errorUtils.EntityError)) {
	u.createUserService = f
}

func (u *userServiceMock) SetUpdateUser(f func(*domain.User) (*domain.User, errorUtils.EntityError)) {
	u.updateUserService = f
}

func (u *userServiceMock) SetDelete(f func(uint64) errorUtils.EntityError) {
	u.deleteUserService = f
}

func (u *userServiceMock) SetGetAll(f func() ([]domain.User, errorUtils.EntityError)) {
	u.getAllUserService = f
}

type UserControllerTestSuite struct {
	suite.Suite
	mockService UserServiceMockInterface
	r *gin.Engine
	rr *httptest.ResponseRecorder

}

func TestUsersControllerTestSuite(t *testing.T){
	suite.Run(t, new(UserControllerTestSuite))
}

func (s *UserControllerTestSuite) SetupSuite() {
	mock := &userServiceMock{}
	s.mockService = mock
	services.UsersService = mock
	s.r = gin.Default()
	router.InitAllUserRoutes(s.r)
}

func (s *UserControllerTestSuite) BeforeTest(_, _ string) {
	s.rr = httptest.NewRecorder()
}

func (s *UserControllerTestSuite) TestGetUser_Success() {
	s.mockService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})
	userIdParam := "1"
	req, _ := http.NewRequest(http.MethodGet, "/users/"+userIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	var user domain.User
	err := json.Unmarshal(s.rr.Body.Bytes(), &user)
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, http.StatusOK, s.rr.Code)
	assert.EqualValues(t, 1, user.ID)
	assert.EqualValues(t, "dev", user.Name)
	assert.EqualValues(t, "dev@test.com", user.Email)
}

func (s *UserControllerTestSuite) TestGetUser_InvalidId() {
	userIdParam := "abc"
	req, _ := http.NewRequest(http.MethodGet, "/users/"+userIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "user id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func (s *UserControllerTestSuite) TestGetUser_NotFound() {
	s.mockService.SetGetUser(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return nil, errorUtils.NewNotFoundError("user not found")
	})
	userIdParam := "1"
	req, _ := http.NewRequest(http.MethodGet, "/users/"+userIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusNotFound, apiErr.Status())
	assert.EqualValues(t, "user not found", apiErr.Message())
	assert.EqualValues(t, "not_found", apiErr.Error())
}

func (s *UserControllerTestSuite) TestGetUser_DatabaseError() {
	s.mockService.SetGetUser(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("database error")
	})
	userIdParam := "1"
	req, _ := http.NewRequest(http.MethodGet, "/users/"+userIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "database error", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
}

func (s *UserControllerTestSuite) TestCreateUser_Success() {
	s.mockService.SetCreateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})
	jsonBody := `{"name":"dev", "email":"dev@test.com"}`
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	var user domain.User
	err = json.Unmarshal(s.rr.Body.Bytes(), &user)
	t:= s.T()
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, http.StatusCreated, s.rr.Code)
	assert.EqualValues(t, 1, user.ID)
	assert.EqualValues(t, "dev", user.Name)
	assert.EqualValues(t, "dev@test.com", user.Email)
}

func (s *UserControllerTestSuite) TestCreateUser_InvalidJsonBadFieldType() {
	jsonBody := `{"name":123456, "email":"dev@test.com"}`
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)
	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func (s *UserControllerTestSuite) TestCreateUser_InvalidJsonMissingField() {
	jsonBody := `{"nam":"dev", "email":"dev@test.com"}`
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)
	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func (s *UserControllerTestSuite) TestUpdateUser_Success() {
	s.mockService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID: 1,
			Name: "dev updated",
			Email: "dev.updated@test.com",
		}, nil
	})
	id := "1"
	jsonBody := `{"name":"dev updated", "email":"dev.updated@test.com"}`
	req, err := http.NewRequest(http.MethodPatch, "/users/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	var user domain.User
	err = json.Unmarshal(s.rr.Body.Bytes(), &user)
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, http.StatusOK, s.rr.Code)
	assert.EqualValues(t, 1, user.ID)
	assert.EqualValues(t, "dev updated", user.Name)
	assert.EqualValues(t, "dev.updated@test.com", user.Email)
}

func (s *UserControllerTestSuite) TestUpdateUser_InvalidId() {
	userIdParam := "abc"
	req, _ := http.NewRequest(http.MethodPatch, "/users/"+userIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "user id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func (s *UserControllerTestSuite) TestUpdateUser_InvalidJson() {
	jsonBody := `{"name":123456, "email":"dev@test.com"}`
	id:="1"
	req, err := http.NewRequest(http.MethodPatch, "/users/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)
	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
	assert.EqualValues(t, "invalid_request", apiErr.Error())
}

func (s *UserControllerTestSuite) TestUpdateUser_ErrorUpdating() {
	s.mockService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error updating user")
	})

	id := "1"
	jsonBody := `{"name":"dev updated", "email":"dev.updated@test.com"}`
	req, err := http.NewRequest(http.MethodPatch, "/users/"+id, bytes.NewBufferString(jsonBody))
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error updating user", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}

func (s *UserControllerTestSuite) TestDeleteUser_Success() {
	s.mockService.SetDelete(func(u uint64) errorUtils.EntityError {
		return nil
	})
	id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/users/"+id, nil)
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	var response = make(map[string]string)
	theErr := json.Unmarshal(s.rr.Body.Bytes(), &response)
	if theErr != nil {
		s.T().Errorf("could not unmarshal response: %v\n", theErr)
	}

	assert.EqualValues(s.T(), http.StatusOK, s.rr.Code)
	assert.EqualValues(s.T(), "deleted", response["status"])
}

func (s *UserControllerTestSuite) TestDeleteUser_InvalidId() {
	userIdParam := "abc"
	req, _ := http.NewRequest(http.MethodDelete, "/users/"+userIdParam, nil)
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "user id should be a number", apiErr.Message())
	assert.EqualValues(t, "bad_request", apiErr.Error())
}

func (s *UserControllerTestSuite) TestDeleteUser_Failure() {
	s.mockService.SetDelete(func(u uint64) errorUtils.EntityError {
		return errorUtils.NewInternalServerError("error deleting user")
	})
	id:="1"
	req, err := http.NewRequest(http.MethodDelete, "/users/"+id, nil)
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error deleting user", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}

func (s *UserControllerTestSuite) TestGetAllUsers_Success() {
	s.mockService.SetGetAll(func() ([]domain.User, errorUtils.EntityError) {
		return []domain.User{
			{
				ID: 1,
				Name: "dev1",
				Email: "dev1@test.com",
			},
			{
				ID: 2,
				Name: "dev2",
				Email: "dev2@test.com",
			},
		}, nil
	})

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	var users []domain.User
	theErr := json.Unmarshal(s.rr.Body.Bytes(), &users)
	if theErr != nil {
		s.T().Errorf("could not unmarshal response: %v\n", theErr)
	}
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.EqualValues(t, users[0].ID, 1)
	assert.EqualValues(t, users[0].Name, "dev1")
	assert.EqualValues(t, users[0].Email, "dev1@test.com")
	assert.EqualValues(t, users[1].ID, 2)
	assert.EqualValues(t, users[1].Name, "dev2")
	assert.EqualValues(t, users[1].Email, "dev2@test.com")
}



func (s *UserControllerTestSuite) TestGetAllUsers_Failure() {
	s.mockService.SetGetAll(func() ([]domain.User, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error getting users")
	})
	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		s.T().Errorf("error while creating the request: %v\n", err)
	}
	s.r.ServeHTTP(s.rr, req)

	apiErr, err := errorUtils.NewApiErrFromBytes(s.rr.Body.Bytes())
	t:=s.T()
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, "error getting users", apiErr.Message())
	assert.EqualValues(t, "server_error", apiErr.Error())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
}
