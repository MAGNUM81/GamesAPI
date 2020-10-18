package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

var (
	tm = time.Now()

)

type UserRepoMockInterface interface {
	SetGetUserDomain(func(uint64) (*domain.User, errorUtils.EntityError))
	SetCreateUserDomain(func(user *domain.User) (*domain.User, errorUtils.EntityError))
	SetUpdateUserDomain(func(user *domain.User) (*domain.User, errorUtils.EntityError))
	SetDeleteUserDomain(func(id uint64) errorUtils.EntityError)
	SetGetAllUserDomain(func() ([]domain.User, errorUtils.EntityError))
}

type userRepoMock struct {
	getUserDomain func(id uint64) (*domain.User, errorUtils.EntityError)
	createUserDomain func(user *domain.User) (*domain.User, errorUtils.EntityError)
	updateUserDomain func(user *domain.User) (*domain.User, errorUtils.EntityError)
	deleteUserDomain func(id uint64) errorUtils.EntityError
	getAllUsersDomain func() ([]domain.User, errorUtils.EntityError)
}

//UserRepoMockInterface implementation, so we can swap the methods around and get the desired behavior from the repository
func (m *userRepoMock) SetGetUserDomain(f func(uint64) (*domain.User, errorUtils.EntityError)) {
	m.getUserDomain = f
}

func (m *userRepoMock) SetCreateUserDomain(f func(user *domain.User) (*domain.User, errorUtils.EntityError)) {
	m.createUserDomain = f
}

func (m *userRepoMock) SetUpdateUserDomain(f func(user *domain.User) (*domain.User, errorUtils.EntityError)) {
	m.updateUserDomain = f
}

func (m *userRepoMock) SetDeleteUserDomain(f func(id uint64) errorUtils.EntityError) {
	m.deleteUserDomain = f
}

func (m *userRepoMock) SetGetAllUserDomain(f func() ([]domain.User, errorUtils.EntityError)) {
	m.getAllUsersDomain = f
}

//UserRepoInterface implementation (redirects all calls to the swappable methods)
func (m *userRepoMock) Get(id uint64) (*domain.User, errorUtils.EntityError){
	return m.getUserDomain(id)
}
func (m *userRepoMock) Create(msg *domain.User) (*domain.User, errorUtils.EntityError){
	return m.createUserDomain(msg)
}
func (m *userRepoMock) Update(msg *domain.User) (*domain.User, errorUtils.EntityError){
	return m.updateUserDomain(msg)
}
func (m *userRepoMock) Delete(id uint64) errorUtils.EntityError {
	return m.deleteUserDomain(id)
}
func (m *userRepoMock) GetAll() ([]domain.User, errorUtils.EntityError) {
	return m.getAllUsersDomain()
}
func (m *userRepoMock) Initialize(_ *gorm.DB) {}

type UserServiceTestSuite struct {
	suite.Suite
	mockRepository UserRepoMockInterface
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) SetupSuite() {
	mock := &userRepoMock{}

	s.mockRepository = mock //set this so we can swap the methods
	domain.UserRepo = mock  //set this so the tested code calls the swapped methods
}

func (s *UserServiceTestSuite) TestUsersService_GetUser_Success() {
	s.mockRepository.SetGetUserDomain(func(userId uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:        1,
			Email:     "devgolang@test.com",
			Name:      "devgolang",
			CreatedAt: tm,
		}, nil
	})
	user, err := services.UsersService.GetUser(1)
	assert.NotNil(s.T(), user)
	assert.Nil(s.T(), err)
	assert.EqualValues(s.T(), 1, user.ID)
	assert.EqualValues(s.T(), "devgolang@test.com", user.Email)
	assert.EqualValues(s.T(), "devgolang", user.Name)
	assert.EqualValues(s.T(), tm, user.CreatedAt)
}

func (s *UserServiceTestSuite) TestUsersService_GetUser_NotFound() {
	expectedError := errorUtils.NewNotFoundError("user was not found")
	s.mockRepository.SetGetUserDomain(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return nil, expectedError
	})
	user, err := services.UsersService.GetUser(1)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), user)
	assert.Equal(s.T(), expectedError, err)
}

func (s *UserServiceTestSuite) TestUsersService_CreateUser_Success() {
	expectedUser := &domain.User{
		ID: 1,
		Email: "dev@test.com",
		Name: "dev",
		CreatedAt: tm,
	}
	s.mockRepository.SetCreateUserDomain(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return expectedUser, nil
	})
	request := &domain.User{
		Email: "dev@test.com",
		Name: "dev",
		CreatedAt: tm,
	}
	user, err := services.UsersService.CreateUser(request)
	assert.NotNil(s.T(), user)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedUser, user)
}

func (s *UserServiceTestSuite) TestUsersService_CreateUser_InvalidRequest() {
	tests := []struct {
		request *domain.User
		expectedError errorUtils.EntityError
	}{
		{
			request: &domain.User{
				Name:     "",
				Email:      "dev@test.com",
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("User name cannot be empty"),
		},
		{
			request: &domain.User{
				Name:     "dev",
				Email:      "",
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("User email cannot be empty"),
		},
		{
			request: &domain.User{
				Name:     "dev",
				Email:      "badly_formatted_email",
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("User email is not formatted correctly"),
		},
	}
	for _, tt := range tests {
		msg, err := services.UsersService.CreateUser(tt.request)
		assert.Nil(s.T(), msg)
		assert.NotNil(s.T(), err)
		assert.Equal(s.T(), tt.expectedError, err)
	}
}

//one reason why the create could fail is a violation of unique constraint on email
//therefore that's what we "test" here.
func (s *UserServiceTestSuite) TestUsersService_CreateUser_Failure() {
	expectedErr := errorUtils.NewInternalServerError("email is already in use")
	s.mockRepository.SetCreateUserDomain(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return nil, expectedErr
	})
	request := &domain.User{
		ID:    1,
		Name:  "dev",
		Email: "dev@test.com",
		CreatedAt: tm,
	}
	user, err := services.UsersService.CreateUser(request)
	assert.Nil(s.T(), user)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedErr, err)
}

func (s *UserServiceTestSuite) TestUsersService_UpdateUser_Success() {
	before := &domain.User{
		ID:    1,
		Name:  "devBefore",
		Email: "devBefore@test.com",
	}
	expectedAfter := &domain.User{
		ID:    1,
		Name:  "devAfter",
		Email: "devAfter@test.com",
	}
	s.mockRepository.SetGetUserDomain(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return before, nil
	})
	s.mockRepository.SetUpdateUserDomain(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return expectedAfter, nil
	})

	request := expectedAfter

	user, err := services.UsersService.UpdateUser(request)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), expectedAfter, user)
}

func (s *UserServiceTestSuite) TestUsersService_UpdateUser_InvalidRequest() {
	tests := []struct {
		request *domain.User
		expectedError errorUtils.EntityError
	}{
		{
			request: &domain.User{
				Name:     "",
				Email:      "dev@test.com",
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("User name cannot be empty"),
		},
		{
			request: &domain.User{
				Name:     "dev",
				Email:      "",
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("User email cannot be empty"),
		},
		{
			request: &domain.User{
				Name:     "dev",
				Email:      "badly_formatted_email",
				CreatedAt: tm,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("User email is not formatted correctly"),
		},
	}
	for _, tt := range tests {
		msg, err := services.UsersService.UpdateUser(tt.request)
		assert.Nil(s.T(), msg)
		assert.NotNil(s.T(), err)
		assert.Equal(s.T(), tt.expectedError, err)
	}
}

func (s *UserServiceTestSuite) TestUsersService_UpdateUser_FailureGettingFormerUser() {
	s.mockRepository.SetGetUserDomain(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error getting user")
	})
	request := &domain.User{
		Name:  "dev",
		Email: "dev@test.com",
	}
	msg, err := services.UsersService.UpdateUser(request)
	t := s.T()
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error getting user", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func (s *UserServiceTestSuite) TestUsersService_UpdateUser_FailureUpdatingUser() {
	s.mockRepository.SetGetUserDomain(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})
	s.mockRepository.SetUpdateUserDomain(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error updating user")
	})

	request := &domain.User{
		ID:    1,
		Name:  "devAAA",
		Email: "devAAA@test.com",
	}
	msg, err := services.UsersService.UpdateUser(request)
	t:= s.T()
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error updating user", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func (s *UserServiceTestSuite) TestUsersService_DeleteUser_Success() {
	s.mockRepository.SetGetUserDomain(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})
	s.mockRepository.SetDeleteUserDomain(func(_ uint64) errorUtils.EntityError {
		return nil
	})

	err := services.UsersService.DeleteUser(1)
	assert.Nil(s.T(), err)
}

func (s *UserServiceTestSuite) TestUsersService_DeleteUser_ErrorGettingUser() {
	expectedError := errorUtils.NewInternalServerError("Something went wrong fetching user")
	s.mockRepository.SetGetUserDomain(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return nil, expectedError
	})
	err := services.UsersService.DeleteUser(1)
	t := s.T()
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func (s *UserServiceTestSuite) TestUsersService_DeleteUser_ErrorDeletingUser(){
	expectedError := errorUtils.NewInternalServerError("error deleting message")
	s.mockRepository.SetGetUserDomain(func(u uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})
	s.mockRepository.SetDeleteUserDomain(func(id uint64) errorUtils.EntityError {
		return expectedError
	})

	err := services.UsersService.DeleteUser(1)
	t := s.T()
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func (s *UserServiceTestSuite) TestUsersService_GetAll_Success() {
	s.mockRepository.SetGetAllUserDomain(func() ([]domain.User, errorUtils.EntityError) {
		return []domain.User{
			{
				ID: 1,
				Email: "dev1@test.com",
				Name: "dev1",
			},
			{
				ID: 2,
				Email: "dev2@test.com",
				Name: "dev2",
			},
			{
				ID: 3,
				Email: "dev3@test.com",
				Name: "dev3",
			},
		}, nil
	})
	users, err := services.UsersService.GetAllUsers()
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.EqualValues(t, users[0].ID, 1)
	assert.EqualValues(t, users[0].Name, "dev1")
	assert.EqualValues(t, users[0].Email, "dev1@test.com")
	assert.EqualValues(t, users[1].ID, 2)
	assert.EqualValues(t, users[1].Name, "dev2")
	assert.EqualValues(t, users[1].Email, "dev2@test.com")
	assert.EqualValues(t, users[2].ID, 3)
	assert.EqualValues(t, users[2].Name, "dev3")
	assert.EqualValues(t, users[2].Email, "dev3@test.com")
}

func (s *UserServiceTestSuite) TestUsersService_GetAllUsers_ErrorGettingUsers() {
	expectedErr := errorUtils.NewInternalServerError("error getting users")
	s.mockRepository.SetGetAllUserDomain(func() ([]domain.User, errorUtils.EntityError) {
		return nil, expectedErr
	})

	users, err := services.UsersService.GetAllUsers()
	t := s.T()
	assert.NotNil(t, err)
	assert.Nil(t, users)
	assert.Equal(t, expectedErr, err)
}
