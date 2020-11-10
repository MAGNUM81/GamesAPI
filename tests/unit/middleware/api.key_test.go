package middleware

import (
	"GamesAPI/src/middleware"
	"GamesAPI/src/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)
type TokenServiceMockInterface interface {
	SetValidateToken(func(string)(bool, error))
}
type tokenServiceMock struct {
	validateToken func(string)(bool, error)
}

func (t tokenServiceMock) SetValidateToken(f func(string) (bool, error)) {
	t.validateToken = f
}

func (t tokenServiceMock) GetApiToken() (token string, err error) {
	return  "1234",nil
}

func (t tokenServiceMock) ValidateToken(s string) (token bool, err error) {
	return t.validateToken(s)
}

type TestMiddlewareApiTokenSuite struct {
	suite.Suite
	mockService TokenServiceMockInterface
	r *gin.Engine
	rr *httptest.ResponseRecorder
}

func (t *TestMiddlewareApiTokenSuite) BeforeTest(_, _ string) {
	t.rr = httptest.NewRecorder()
}

func TestMiddlewareApiTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TestMiddlewareApiTokenSuite))
}

func (t *TestMiddlewareApiTokenSuite) SetupSuite() {
	mock := &tokenServiceMock{}

	t.mockService = mock //set this so we can swap the methods
	services.TokenService = mock  //set this so the tested code calls the swapped methods
	t.r = gin.Default()
	t.r.Use(middleware.MiddlewareHandler)

}


//Test identitfaction Token
//this test is a correct api token
func (t *TestMiddlewareApiTokenSuite) TestMiddlewareService_Authentication_TokenValid(){
	t.mockService.SetValidateToken(func(string) (bool, error){
		return  true , nil

	})

	//1. api token - n'importe quelle valeur string
	 tokenApi := "1234"
	//2. nouvelle requête http semblable à http.NewRequest(http.MethodGet, "/", nil)
	//2.1 mettre api token dans header de la requête - google it
	req, _ :=http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("api_token",tokenApi)
	t.r.ServeHTTP(t.rr,req)
	//3. "s'envoyer" la requête - semblable à s.r.ServeHTTP(s.rr, req)
	//3.1 s.r : *gin.Engine
	//3.2 s.rr : *httptest.ResponseRecorder
	//3.3 ne pas oublier de les initialiser
	//SetupSuite -> s.r = gin.Default()
	//3.4 r.GET("/", middleware.MiddlewareHandler)
	//BeforeTest -> s.rr = httptest.NewRecorder()



}