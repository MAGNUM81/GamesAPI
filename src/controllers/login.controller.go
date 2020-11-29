package controllers

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func abortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"Error": message})
}

func LoginController(c *gin.Context)  {
	authenticateHeader := strings.Trim(c.Request.Header.Get("Www-Authenticate"), " ")
	parts := strings.Split(authenticateHeader, ";")
	if authenticateHeader == "" || len(parts) != 2 {
		abortWithError(c, http.StatusBadRequest, "Www-Authenticate header was not set properly")
	}

	email, password := parts[0], []byte(parts[1])

	//try to find user with email

	users, err := services.UsersService.GetAllUsers()

	if err != nil {
		abortWithError(c, err.Status(), err.Message())
		return
	}

	var potentialUser *domain.User = nil

	for _, user := range users {
		if user.Email == email {
			potentialUser = &user
			break
		}
	}

	if potentialUser == nil {
		//no user has been found for email
		abortWithError(c, http.StatusUnauthorized, "Bad username/password combination")
		return
	}

	validPasswordHash := potentialUser.PasswordHash

	isPasswordValid, authErr := services.AuthenticationService.ValidatePassword(password, validPasswordHash)

	if authErr != nil {
		//password doesn't match
		abortWithError(c, http.StatusUnauthorized, "Bad username/password combination")
		return
	}

	if !isPasswordValid {
		//typically this shouldn't happen, but we put it here just in case
		abortWithError(c, http.StatusUnauthorized, "Bad username/password combination")
		return
	}

	//To be documented : a typical session lasts 10 minutes.
	//					 if a refresh is issued after the 10-minute mark, the user will have to authenticate again.
	expireAt := time.Now().Add(time.Minute * 10)
	token, tokenErr := services.UserSessionService.GenerateSessionToken(potentialUser.ID, expireAt)

	if tokenErr != nil {
		abortWithError(c, http.StatusInternalServerError, fmt.Sprintf("Token couldn't be generated by server - %s", tokenErr.Error()))
		return
	}

	session := &domain.UserSession{
		Token:     token,
		UserId:    potentialUser.ID,
		ExpiresAt: expireAt.UnixNano(),
	}

	_, createSessionError := services.UserSessionService.CreateSession(session)

	//TODO: delete older session if present, to prevent same user from having many session tokens at a time.

	if createSessionError != nil {
		abortWithError(c, createSessionError.Status(), createSessionError.Message())
		return
	}

	c.Header("Authorization", token)

	c.JSON(http.StatusOK, gin.H{"Message" : fmt.Sprintf("User with email '%s' Successfully authenticated. " +
		"Session token was sent in response's 'Authorization' header ", potentialUser.Email)})
}