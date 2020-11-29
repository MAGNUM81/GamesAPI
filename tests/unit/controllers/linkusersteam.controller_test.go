package controllers

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/router"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

///func TestUsersControllerTestSuite(t *testing.T){
//	suite.Run(t, new(UserControllerTestSuite))
//}

//func (s *UserControllerTestSuite) SetupSuite() {
	//mock := &userServiceMock{}
	//s.mockService = mock
	//services.UsersService = mock
	///s.r = gin.Default()
	//router.InitAllUserRoutes(s.r)
//}

func (s *UserControllerTestSuite) TestLinkUserSteam_ValidSteamId() {
	s.mockService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

		usersteamid :="12345678911234567"
		req, _ := http.NewRequest(http.MethodGet, "/", "{\"profile_url\":\" https://steamcommunity.com/profiles/12345678911234567\"}")
		s.r.ServeHTTP(s.rr, req)

	//assert.EqualValues(usersteamid,)

}