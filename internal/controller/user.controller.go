package controller

import (
	"Gober/common/response"
	"Gober/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (uc *UserController) GetUserById(c *gin.Context) {
	id := c.Query("id")
	
	response.SuccessResponse(c, response.ErrCodeSuccess, uc.userService.GetUserById(id))

}