package middleware

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/middleware"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"GamesAPI/tests/unit/mocks"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type AuthTestSuite struct {
	suite.Suite
	mockAuthService     mocks.AuthorizationServiceMockInterface
	mockUserRoleService mocks.UserRoleServiceMockInterface
	r                   *gin.Engine
	rr                  *httptest.ResponseRecorder
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupSuite() {
	mockAuth := &mocks.AuthorizationServiceMock{}
	mockUserRoles := &mocks.UserRoleMock{}
	s.mockAuthService = mockAuth
	services.AuthorizationService = mockAuth
	s.mockUserRoleService = mockUserRoles
	services.UserRoleService = mockUserRoles
	s.r = gin.Default()
	s.r.Use(middleware.AuthorizationHandler)
	s.r.GET("/games", BidonController)
	s.r.GET("/achievements", BidonController)
	s.r.HEAD("/games", BidonController)

}

func BidonController(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ok"})
}

func (s *AuthTestSuite) BeforeTest(_, _ string) {
	s.rr = httptest.NewRecorder()
}

func (s *AuthTestSuite) TestAuth_UserNotExists() {
	s.mockUserRoleService.SetGetRolesByUserID(func(userId uint64) ([]domain.UserRole, errorUtils.EntityError) {
		return nil, errorUtils.NewNotFoundError("user cannot be found")
	})
	req, _ := http.NewRequest(http.MethodGet, "/games", nil)
	req = req.WithContext(context.WithValue(context.Background(), domain.RbacUserId(), uint64(1)))
	s.r.ServeHTTP(s.rr, req)

	t := s.T()
	assert.EqualValues(t, 404, s.rr.Code)

}

func (s *AuthTestSuite) TestAuth_NoRoles() {
	s.mockUserRoleService.SetGetRolesByUserID(func(userId uint64) ([]domain.UserRole, errorUtils.EntityError) {
		return []domain.UserRole{}, nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/games", nil)
	req = req.WithContext(context.WithValue(context.Background(), domain.RbacUserId(), uint64(1)))
	s.r.ServeHTTP(s.rr, req)

	t := s.T()
	assert.EqualValues(t, 500, s.rr.Code)
}

func (s *AuthTestSuite) TestAuth_BadResource() {
	s.mockUserRoleService.SetGetRolesByUserID(func(userId uint64) ([]domain.UserRole, errorUtils.EntityError) {
		return []domain.UserRole{
			{
				ID:     1,
				UserID: 1,
				Name:   "Admin",
			},
		}, nil
	})
	resource := "/achievements"
	req, _ := http.NewRequest(http.MethodGet, resource, nil)
	req = req.WithContext(context.WithValue(context.Background(), domain.RbacUserId(), uint64(1)))
	s.r.ServeHTTP(s.rr, req)

	t := s.T()
	assert.EqualValues(t, 400, s.rr.Code)
}

func (s *AuthTestSuite) TestAuth_BadEndpoint() {
	s.mockUserRoleService.SetGetRolesByUserID(func(userId uint64) ([]domain.UserRole, errorUtils.EntityError) {
		return []domain.UserRole{
			{
				ID:     1,
				UserID: 1,
				Name:   "Admin",
			},
		}, nil
	})
	endpoint := http.MethodHead
	req, _ := http.NewRequest(endpoint, "/games", nil)
	req = req.WithContext(context.WithValue(context.Background(), domain.RbacUserId(), uint64(1)))
	s.r.ServeHTTP(s.rr, req)

	t := s.T()
	assert.EqualValues(t, 400, s.rr.Code)
}

func (s *AuthTestSuite) TestAuth_ForbiddenAccess() {
	s.mockUserRoleService.SetGetRolesByUserID(func(userId uint64) ([]domain.UserRole, errorUtils.EntityError) {
		return []domain.UserRole{
			{
				ID:     1,
				UserID: 1,
				Name:   "User",
			},
		}, nil
	})
	s.mockAuthService.SetAuthorize(func(ctx context.Context, url *url.URL, role string, resource string, endpoint string) error {
		return errors.New("query does not comply")
	})
	resource := "/games"
	endpoint := http.MethodGet

	req, _ := http.NewRequest(endpoint, resource, nil)
	req = req.WithContext(context.WithValue(context.Background(), domain.RbacUserId(), uint64(1)))
	s.r.ServeHTTP(s.rr, req)

	t := s.T()
	assert.EqualValues(t, 403, s.rr.Code)
}

func (s *AuthTestSuite) TestAuth_GrantedAccess() {
	s.mockUserRoleService.SetGetRolesByUserID(func(userId uint64) ([]domain.UserRole, errorUtils.EntityError) {
		return []domain.UserRole{
			{
				ID:     1,
				UserID: 1,
				Name:   "Admin",
			},
		}, nil
	})
	s.mockAuthService.SetAuthorize(func(ctx context.Context, url *url.URL, role string, resource string, endpoint string) error {
		return nil
	})
	resource := "/games"
	endpoint := http.MethodGet

	req, _ := http.NewRequest(endpoint, resource, nil)
	req = req.WithContext(context.WithValue(context.Background(), domain.RbacUserId(), uint64(1)))
	s.r.ServeHTTP(s.rr, req)

	t := s.T()
	assert.EqualValues(t, 200, s.rr.Code)
}
