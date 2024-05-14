package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	UserService
}

func NewController(s UserService) *Controller {
	return &Controller{
		UserService: s,
	}
}

func (controller *Controller) CreateUser(c *gin.Context) {
	var createUserReq CreateUserReq
	if err := c.ShouldBindJSON(&createUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := controller.UserService.CreateUser(c.Request.Context(), &createUserReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (controller *Controller) Login(c *gin.Context) {
	var loginUserReq LoginUserReq
	if err := c.ShouldBindJSON(&loginUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := controller.UserService.Login(c.Request.Context(), &loginUserReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", res.AccessToken, 60*60*24, "/", "localhost", false, true) // Secure = false, httpOnly = true
	c.JSON(http.StatusOK, res)
}

func (controller *Controller) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
