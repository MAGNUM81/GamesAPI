package controllers

import (
	"GamesAPI/src/External/Steam"
	"GamesAPI/src/controllers"
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"GamesAPI/tests/unit/mocks"
	"bytes"
	"errors"
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
	mockUserService mocks.UserServiceMockInterface
	mockSteamService mocks.SteamUserMockInterface
	r *gin.Engine
	rr *httptest.ResponseRecorder
}

func TestLinkSteamUsersControllerTestSuite(t *testing.T){
	suite.Run(t, new(LinkSteamUserTestSuite))
}

func (s *LinkSteamUserTestSuite) SetupSuite() {
	mock := &mocks.UserServiceMock{}
	steammock :=&mocks.SteamUserMock{}
	s.mockSteamService = steammock
	s.mockUserService = mock
	Steam.ExternalSteamUserService = steammock
	services.UsersService = mock
	s.r = gin.Default()
	s.r.POST("/", controllers.LinkSteamUser)
}

func (s *LinkSteamUserTestSuite) BeforeTest(_, _ string) {
	s.rr = httptest.NewRecorder()
}

func (s *LinkSteamUserTestSuite) TestLinkUserSteam_ValidProfile() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	usersteamid := "12345678911234567"
	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://steamcommunity.com/profiles/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 200, s.rr.Code)

}

func (s *LinkSteamUserTestSuite) TestLinkUserSteam_InvalidProfile() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 400, s.rr.Code)

}

func (s *LinkSteamUserTestSuite) TestLinkUserSteam_ValidSteamId() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	usersteamid := "gabelogannewell"
	s.mockSteamService.SetGetUserID(func(usersteamid string) (string, error){
		return "12345678911234567", nil
	})


	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://steamcommunity.com/id/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 200, s.rr.Code)}

func (s *LinkSteamUserTestSuite) TestLinkUserSteam_invalidSteamId() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	usersteamid := "gabelogannewell"
	s.mockSteamService.SetGetUserID(func(usersteamid string) (string, error){
		return "", errors.New("error steam id invalid")
	})


	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://steamcommunity.com/id/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 500 , s.rr.Code)}

// steam url  component id or profile is empty
func (s *LinkSteamUserTestSuite) TestLinkUserSteam_InvalidValidateUrl_null() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	usersteamid := ""
	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://steamcommunity.com/profiles/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 400, s.rr.Code)
}

// steam url steamcommunity change to unrealcommunity
func (s *LinkSteamUserTestSuite) TestLinkUserSteam_InvalidValidateUrl_steamcommunity() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	usersteamid := "12345678911234567"

	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://unrealcommunity.com/profiles/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 400, s.rr.Code)

}

// steam url Number of components is incorrect normal len = 5 this 4
func (s *LinkSteamUserTestSuite) TestLinkUserSteam_InvalidValidateUrl_len() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	usersteamid := "12345678911234567"

	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https:/steamcommunity.com/profiles/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 400, s.rr.Code)

}
// steam url id or profile is not the 4 componant
func (s *LinkSteamUserTestSuite) TestLinkUserSteam_InvalidValidateUrl_component_notfound() {
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})

	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	usersteamid := "12345678911234567"

	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://steamcommunity.com/whatsup/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 400, s.rr.Code)

}



func (s *LinkSteamUserTestSuite) TestLinkUserSteam_InvalidGetUserid() {
	expectedErr := errorUtils.NewInternalServerError("error getting users")
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return nil, expectedErr
	})


	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:   1 ,
			Name:  "dev",
			Email: "dev@test.ru",
		}, nil
	})

	usersteamid := "12345678911234567"

	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://steamcommunity.com/profiles/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 500, s.rr.Code)

}
func (s *LinkSteamUserTestSuite) TestLinkUserSteam_InvalidUpdateUser() {
	expectedErr := errorUtils.NewInternalServerError("error updating users")
	s.mockUserService.SetGetUser(func(id uint64) (*domain.User, errorUtils.EntityError) {
		return &domain.User{
			ID:    1,
			Name:  "dev",
			Email: "dev@test.com",
		}, nil
	})


	s.mockUserService.SetUpdateUser(func(user *domain.User) (*domain.User, errorUtils.EntityError) {
		return nil, expectedErr
	})

	usersteamid := "12345678911234567"

	jsonBody := bytes.NewBufferString(fmt.Sprintf(`{"profile_url":"https://steamcommunity.com/profiles/%s", "userid":1}`, usersteamid))
	req, _ := http.NewRequest(http.MethodPost, "/", jsonBody)
	s.r.ServeHTTP(s.rr, req)

	assert.EqualValues(s.T(), 500, s.rr.Code)

}