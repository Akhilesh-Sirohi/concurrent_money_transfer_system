package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"concurrent_money_transfer_system/utils"
)

type UserController struct {
	userService UserService
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user User
	err := utils.BindAndValidateRequest(c, &user)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	createdUser, err := uc.userService.CreateUser(c.Request.Context(), user)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResponseError(c, utils.NewErrorWithMessage(utils.ErrInvalidRequest, "User ID is required"))
		return
	}

	user, err := uc.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user User
	err := utils.BindAndValidateRequest(c, &user)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	updatedUser, err := uc.userService.UpdateUser(c.Request.Context(), user)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResponseError(c, utils.NewErrorWithMessage(utils.ErrInvalidRequest, "User ID is required"))
		return
	}

	err := uc.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "User deleted successfully"})
}

func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}
