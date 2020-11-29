package controllers

import (
	"GamesAPI/src/controllers"
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"GamesAPI/tests/unit/mocks"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LinkSteamUserTestSuite struct {
	suite.Suite
	mockUserService UserServiceMockInterface
	mockSteamService mocks.SteamUserMockInterface
	r *gin.Engine
	rr *httptest.ResponseRecorder
}

func TestLinkSteamUsersControllerTestSuite(t *testing.T){
	suite.Run(t, new(LinkSteamUserTestSuite))
}

func (s *LinkSteamUserTestSuite) SetupSuite() {
	mock := &userServiceMock{}
	s.mockUserService = mock
	services.UsersService = mock
	s.r = gin.Default()
	s.r.POST("/", controllers.LinkSteamUser)
}

func (s *LinkSteamUserTestSuite) BeforeTest() {
	s.rr = httptest.NewRecorder()
}

func (s *LinkSteamUserTestSuite) TestLinkUserSteam_ValidSteamId() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})
	usersteamid := "12345678911234567"
	req, _ := http.NewRequest(http.MethodGet, "/", bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://steamcommunity.com/profiles/%s"}`, usersteamid)))
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 200, s.rr.Code)

}
