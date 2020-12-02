package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"GamesAPI/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type UserRoleServiceTestSuite struct {
	suite.Suite
	mockRepository mocks.UserRoleRepoMockInterface
}

//5. replace common repo by mock in suite setup
//once this is written, steps 3 and 4 can be auto-generated
func (s *UserRoleServiceTestSuite) SetupSuite() {
	mock := &mocks.UserRoleRepoMock{}
	s.mockRepository = mock
	domain.UserRoleRepo = mock
}

func TestRoleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserRoleServiceTestSuite))
}

//6. do tests

//test data

var (
	testUser = domain.UserRole{
		ID:     1,
		UserID: 1,
		Name:   "User",
	}
	testAdmin = domain.UserRole{
		ID:     2,
		UserID: 2,
		Name:   "Admin",
	}
	testManyReturn = []domain.UserRole{
		testUser,
		testAdmin,
	}
)

func (s *UserRoleServiceTestSuite) TestUserRoleService_GetUserRole_NotFound() {
	expectedError := errorUtils.NewNotFoundError("role was not found")
	s.mockRepository.SetGetRole(func(u uint64) (*domain.UserRole, errorUtils.EntityError) {
		return nil, expectedError
	})
	role, err := services.UserRoleService.GetRole(1)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), role)
	assert.Equal(s.T(), expectedError, err)
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_CreateUserRole_Success() {
	expectedUserRole := &testAdmin
	s.mockRepository.SetCreateRole(func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
		return expectedUserRole, nil
	})
	request := &testAdmin
	role, err := services.UserRoleService.CreateRole(request)
	assert.NotNil(s.T(), role)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedUserRole, role)
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_CreateUserRole_InvalidRequest() {
	tests := []struct {
		request       *domain.UserRole
		expectedError errorUtils.EntityError
	}{
		{
			request: &domain.UserRole{
				ID:     1,
				UserID: 1,
				Name:   "",
			},
			expectedError: errorUtils.NewUnprocessableEntityError("role name cannot be empty"),
		},
	}
	for _, tt := range tests {
		s.mockRepository.SetCreateRole(func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
			return nil, tt.expectedError
		})
		msg, err := services.UserRoleService.CreateRole(tt.request)
		assert.Nil(s.T(), msg)
		assert.NotNil(s.T(), err)
		assert.Equal(s.T(), tt.expectedError, err)
	}
}

//one reason why the create could fail is a violation of unique constraint on User ID
//therefore that's what we "test" here.
func (s *UserRoleServiceTestSuite) TestUserRoleService_CreateUserRole_Failure() {
	expectedErr := errorUtils.NewInternalServerError("This user already has a role")
	s.mockRepository.SetCreateRole(func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
		return nil, expectedErr
	})
	request := &testAdmin
	role, err := services.UserRoleService.CreateRole(request)
	assert.Nil(s.T(), role)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedErr, err)
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_UpdateUserRole_Success() {
	before := &testAdmin
	expectedAfter := &domain.UserRole{
		ID:     2,
		UserID: 2,
		Name:   "User",
	}
	s.mockRepository.SetGetRole(func(u uint64) (*domain.UserRole, errorUtils.EntityError) {
		return before, nil
	})
	s.mockRepository.SetUpdateRole(func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
		return expectedAfter, nil
	})

	request := expectedAfter

	role, err := services.UserRoleService.UpdateRole(request)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), role)
	assert.Equal(s.T(), expectedAfter, role)
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_UpdateUserRole_InvalidRequest() {
	tests := []struct {
		request       *domain.UserRole
		expectedError errorUtils.EntityError
	}{
		{
			request: &domain.UserRole{
				Name:   "",
				UserID: 1,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("UserRole name cannot be empty"),
		},
		{
			request: &domain.UserRole{
				Name:   "Not Admin or User",
				UserID: 3,
			},
			expectedError: errorUtils.NewUnprocessableEntityError("UserRole name must be 'User' or 'Admin'"),
		},
	}
	s.mockRepository.SetGetRole(func(u uint64) (*domain.UserRole, errorUtils.EntityError) {
		return &testUser, nil
	})
	for _, tt := range tests {
		s.mockRepository.SetUpdateRole(func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
			return nil, tt.expectedError
		})
		msg, err := services.UserRoleService.UpdateRole(tt.request)
		assert.Nil(s.T(), msg)
		assert.NotNil(s.T(), err)
		assert.Equal(s.T(), tt.expectedError, err)
	}
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_UpdateUserRole_FailureGettingFormerUserRole() {
	s.mockRepository.SetGetRole(func(u uint64) (*domain.UserRole, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error getting role")
	})
	request := &domain.UserRole{
		Name:   "Admin",
		UserID: 2,
	}
	msg, err := services.UserRoleService.UpdateRole(request)
	t := s.T()
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error getting role", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_UpdateUserRole_FailureUpdatingUserRole() {
	s.mockRepository.SetGetRole(func(u uint64) (*domain.UserRole, errorUtils.EntityError) {
		return &domain.UserRole{
			ID:     1,
			Name:   "Admin",
			UserID: 1,
		}, nil
	})
	s.mockRepository.SetUpdateRole(func(role *domain.UserRole) (*domain.UserRole, errorUtils.EntityError) {
		return nil, errorUtils.NewInternalServerError("error updating role")
	})

	request := &domain.UserRole{
		ID:     1,
		Name:   "User",
		UserID: 1,
	}
	msg, err := services.UserRoleService.UpdateRole(request)
	t := s.T()
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "error updating role", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_DeleteUserRole_Success() {
	s.mockRepository.SetGetRole(func(u uint64) (*domain.UserRole, errorUtils.EntityError) {
		return &domain.UserRole{
			ID:     1,
			Name:   "Admin",
			UserID: 1,
		}, nil
	})
	s.mockRepository.SetDeleteRole(func(_ uint64) errorUtils.EntityError {
		return nil
	})

	err := services.UserRoleService.DeleteRole(1)
	assert.Nil(s.T(), err)
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_DeleteUserRole_ErrorGettingUserRole() {
	expectedError := errorUtils.NewInternalServerError("Something went wrong fetching role")
	s.mockRepository.SetGetRole(func(u uint64) (*domain.UserRole, errorUtils.EntityError) {
		return nil, expectedError
	})
	err := services.UserRoleService.DeleteRole(1)
	t := s.T()
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_DeleteUserRole_ErrorDeletingUserRole() {
	expectedError := errorUtils.NewInternalServerError("error deleting message")
	s.mockRepository.SetGetRole(func(u uint64) (*domain.UserRole, errorUtils.EntityError) {
		return &domain.UserRole{
			ID:     1,
			Name:   "Admin",
			UserID: 1,
		}, nil
	})
	s.mockRepository.SetDeleteRole(func(id uint64) errorUtils.EntityError {
		return expectedError
	})

	err := services.UserRoleService.DeleteRole(1)
	t := s.T()
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_GetAll_Success() {
	s.mockRepository.SetGetAllRoles(func() ([]domain.UserRole, errorUtils.EntityError) {
		return testManyReturn, nil
	})
	roles, err := services.UserRoleService.GetAllRoles()
	t := s.T()
	assert.Nil(t, err)
	assert.NotNil(t, roles)
	assert.EqualValues(t, roles[0].ID, 1)
	assert.EqualValues(t, roles[0].Name, "User")
	assert.EqualValues(t, roles[0].UserID, 1)
	assert.EqualValues(t, roles[1].ID, 2)
	assert.EqualValues(t, roles[1].Name, "Admin")
	assert.EqualValues(t, roles[1].UserID, 2)
}

func (s *UserRoleServiceTestSuite) TestUserRoleService_GetAllUserRoles_ErrorGettingUserRoles() {
	expectedErr := errorUtils.NewInternalServerError("error getting roles")
	s.mockRepository.SetGetAllRoles(func() ([]domain.UserRole, errorUtils.EntityError) {
		return nil, expectedErr
	})

	roles, err := services.UserRoleService.GetAllRoles()
	t := s.T()
	assert.NotNil(t, err)
	assert.Nil(t, roles)
	assert.Equal(t, expectedErr, err)
}
