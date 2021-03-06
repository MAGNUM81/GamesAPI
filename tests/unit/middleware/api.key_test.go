package middleware

import (
	"GamesAPI/src/middleware"
	"GamesAPI/src/services"
	"GamesAPI/tests/unit/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestMiddlewareApiTokenSuite struct {
	suite.Suite
	mockService mocks.TokenServiceMockInterface
	r           *gin.Engine
	rr          *httptest.ResponseRecorder
}

func BidonHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"Error": "None"})
}

func (t *TestMiddlewareApiTokenSuite) BeforeTest(_, _ string) {
	t.rr = httptest.NewRecorder()
}

func TestMiddlewareApiTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TestMiddlewareApiTokenSuite))
}

func (t *TestMiddlewareApiTokenSuite) SetupSuite() {
	mock := &mocks.TokenServiceMock{}

	t.mockService = mock         //set this so we can swap the methods
	services.TokenService = mock //set this so the tested code calls the swapped methods
	t.r = gin.Default()
	t.r.Use(middleware.MiddlewareHandler)
	t.r.GET("/", BidonHandler)

}

func (t *TestMiddlewareApiTokenSuite) TestMiddlewareService_Authentication_TokenValid() {
	t.mockService.SetValidateToken(func(string) (bool, error) {
		return true, nil
	})

	tokenApi := "1234"
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("x-api-key", tokenApi)
	t.r.ServeHTTP(t.rr, req)

	assert.Equal(t.T(), http.StatusOK, t.rr.Code)
	assert.NotEqual(t.T(), 400, t.rr.Code)
	assert.NotEqual(t.T(), 401, t.rr.Code)
	assert.NotEqual(t.T(), 500, t.rr.Code)
}
func (t *TestMiddlewareApiTokenSuite) TestMiddlewareService_Authentication_MissingToken() {
	t.mockService.SetValidateToken(func(string) (bool, error) {
		return false, nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	t.r.ServeHTTP(t.rr, req)

	assert.Equal(t.T(), 400, t.rr.Code)
	assert.NotEqual(t.T(), 500, t.rr.Code)
	assert.NotEqual(t.T(), 401, t.rr.Code)
	assert.NotEqual(t.T(), http.StatusOK, t.rr.Code)
}

func (t *TestMiddlewareApiTokenSuite) TestMiddlewareService_Authentication_ErrorValidationToken() {
	t.mockService.SetValidateToken(func(string) (bool, error) {
		return true, errors.New("Environment  Api key Token is not find.")
	})
	tokenApi := "1245"
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("x-api-key", tokenApi)
	t.r.ServeHTTP(t.rr, req)

	assert.Equal(t.T(), 500, t.rr.Code)
	assert.NotEqual(t.T(), 400, t.rr.Code)
	assert.NotEqual(t.T(), 401, t.rr.Code)
	assert.NotEqual(t.T(), http.StatusOK, t.rr.Code)

}

func (t *TestMiddlewareApiTokenSuite) TestMiddlewareService_Authentication_InvalidToken() {
	t.mockService.SetValidateToken(func(string) (bool, error) {
		return false, nil
	})
	tokenApi := "1245"
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("x-api-key", tokenApi)
	t.r.ServeHTTP(t.rr, req)

	assert.Equal(t.T(), 401, t.rr.Code)
	assert.NotEqual(t.T(), 400, t.rr.Code)
	assert.NotEqual(t.T(), 500, t.rr.Code)
	assert.NotEqual(t.T(), http.StatusOK, t.rr.Code)

}
