package controllers

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
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
	var user domain.User
	if err := c.ShouldBindJSON(&user); user.Validate() != nil || err != nil {
		userErr := errorUtils.NewUnprocessableEntityError("invalid json body")
		c.JSON(userErr.Status(), userErr)
		return
	}

	u, err := services.UsersService.CreateUser(&user)
	if errorUtils.IsEntityError(c, err) {
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
