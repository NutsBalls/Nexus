package controllers

import (
	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UserController handles user-related API requests
type UserController struct {
	UserService services.UserService
}

// Register godoc
// @Summary Создает нового пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/register [post]
func (uc *UserController) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := uc.UserService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

// Login godoc
// @Summary Аутентифицирует пользователя и возвращает JWT
// @Description Логин пользователя и получение JWT токена для аутентификации
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/login [post]
func (uc *UserController) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	token, err := uc.UserService.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
