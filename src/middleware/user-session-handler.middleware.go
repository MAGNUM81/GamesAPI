package middleware

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitUserSessionHandler(g ...*gin.RouterGroup) {
	for _, group := range g {
		group.Use(UserSessionHandler)
	}
}

func UserSessionHandler(c *gin.Context) {
	authHeader := c.Request.Header["Authorization"]
	if len(authHeader) < 1 {
		err := errorUtils.NewBadRequestError("Authorization header was not set properly.")
		AbortWithError(c, err.Status(), err.Message())
		return
	}
	reqToken := strings.Trim(c.Request.Header["Authorization"][0], " ")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		_ = fmt.Sprintf("%v", splitToken)
		err := errorUtils.NewBadRequestError("Authorization header was not set properly.")
		AbortWithError(c, err.Status(), err.Message())
		return
	}
	sessionKey := splitToken[1]
	if sessionKey == "" {
		err := errorUtils.NewBadRequestError("Authorization header was not set properly.")
		AbortWithError(c, err.Status(), err.Message())
		return
	}

	if !services.UserSessionService.ExistsSession(sessionKey) {
		AbortWithWWWAuthenticate(c, 401, "session does not exist for given token")
		return
	}

	sessionExpired, err := services.UserSessionService.IsSessionExpired(sessionKey, time.Now())
	if err != nil {
		AbortWithError(c, err.Status(), err.Message())
		return
	}

	if sessionExpired {
		AbortWithWWWAuthenticate(c, 401, "Session is expired")
		return
	}

	session, err := services.UserSessionService.GetSession(sessionKey)
	if err != nil {
		AbortWithWWWAuthenticate(c, 401, "Session doesn't exists, WTF")
		return
	}

	c.Request = c.Request.WithContext(context.WithValue(context.Background(), domain.RbacUserId(), session.UserId))

	c.Next()
}

func AbortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"Error": message})
}

func AbortWithWWWAuthenticate(c *gin.Context, code int, authenticateMessage string) {
	c.Header("Www-ValidatePassword", authenticateMessage)
	c.AbortWithStatusJSON(code, gin.H{"Message": authenticateMessage})
}
