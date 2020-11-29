package controllers

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/authUtils"
	"GamesAPI/src/utils/errorUtils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getUserId(userIdParam string) (uint64, errorUtils.EntityError) {
	userId, userError := strconv.ParseUint(userIdParam, 10, 64)
	if userError != nil {
		return 0, errorUtils.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func GetUser(c *gin.Context) {
	userId, userErr := getUserId(c.Param("id"))
	if errorUtils.IsEntityError(c, userErr) {
		return
	}

	user, err := services.UsersService.GetUser(userId)
	if errorUtils.IsEntityError(c, err) {
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	users, err := services.UsersService.GetAllUsers()
	if errorUtils.IsEntityError(c, err){
		return
	}

	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	body := map[string]string{}
	if err := c.ShouldBindJSON(&body); err != nil {
		userErr := errorUtils.NewUnprocessableEntityError("invalid json body")
		c.JSON(userErr.Status(), userErr)
		return
	}
	if len(body) < 4 {
		c.JSON(http.StatusBadRequest, errorUtils.NewUnprocessableEntityError("invalid json body"))
		return
	}

	/*
	now we should expect a user to be created with this request body :
	{
		"name":"<name>"
		"email" : "<email>"
		"password" : "<password>"
		"role" : "<roleName>"
	}
	*/
	name := body["name"]
	email := body["email"]
	password := body["password"]
	role := body["role"]
	passwordHash, hashErr := authUtils.HashAndSalt([]byte(password))

	if hashErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error" : hashErr.Error()})
		return
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := user.Validate(); err != nil {
		c.JSON(err.Status(), gin.H{"Error" : err.Error()})
		return
	}

	u, err := services.UsersService.CreateUser(user)
	if errorUtils.IsEntityError(c, err) {
		return
	}

	r, err := services.UserRoleService.CreateRole(&domain.UserRole{
		UserID:    user.ID,
		Name:      role,
	})

	if err != nil {
		c.JSON(err.Status(), gin.H{"Error" : err.Error()})
		return
	}

	if r == nil {
		//we should rollback our changes to the user table... the user won't have any role.
		c.JSON(http.StatusInternalServerError, gin.H{"Error":"Could not create role for user"})
		return
	}

	c.JSON(http.StatusCreated, u)
}

func UpdateUser(c *gin.Context) {
	userId, err := getUserId(c.Param("id"))
	if errorUtils.IsEntityError(c, err){
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); user.Validate() != nil || err != nil {
		userErr := errorUtils.NewUnprocessableEntityError("invalid json body")
		c.JSON(userErr.Status(), userErr)
		return
	}
	user.ID = userId
	u, err := services.UsersService.UpdateUser(&user)
	if errorUtils.IsEntityError(c, err) {
		return
	}

	c.JSON(http.StatusOK, u)
}

func DeleteUser(c *gin.Context) {
	userId, err := getUserId(c.Param("id"))
	if errorUtils.IsEntityError(c, err){
		return
	}
	if err := services.UsersService.DeleteUser(userId); errorUtils.IsEntityError(c, err){
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
