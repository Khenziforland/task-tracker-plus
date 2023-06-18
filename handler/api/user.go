package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
	// Create new variable
	var user model.UserLogin

	// Model Binding
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	// Check if Email of Password is empty
	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("email or password is empty"))
		return
	}

	// Create new variable for Model User
	modelUser := &model.User{
		Email: user.Email,
		Password: user.Password,
	}

	// func Login from Service
	token, err := u.userService.Login(modelUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	// Cookie
	cookie, err := c.Cookie("session_token")
	if err != nil {
		cookie = *token
		c.SetCookie("session_token", cookie, 3600, "/", "", false, true)
	}

	// Response Message
	c.JSON(http.StatusOK, model.NewSuccessResponse("login success"))
}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
	userTaskCategory, err := u.userService.GetUserTaskCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, userTaskCategory)
}
