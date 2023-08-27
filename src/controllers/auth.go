package controllers

import (
	apperr "api/src/app_error"
	"api/src/entities"
	"api/src/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	usecases.Authenticator
}

func New(usecase usecases.Authenticator) *AuthController {
	return &AuthController{usecase}
}

// SignUp godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        authBody    body     entities.CreateAuth  true  "Create auth body"
// @Success      201  {object}  entities.AuthEntity
// @Failure      400  {object}  apperr.AppError
// @Failure      500  {object}  apperr.AppError
// @Router       /signup [post]
func (ac *AuthController) SignUp(c echo.Context) error {
	payload := new(entities.CreateAuth)
	if err := c.Bind(&payload); err != nil {
		return apperr.New("100002", http.StatusBadRequest, "Invalid body", "Invalid body", err)
	}

	authEntity, err := ac.Authenticator.SignUp(*payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, authEntity)
}

// SignIn godoc
// @Summary      Authenticate user
// @Description  Decide whether user can login to the system
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login    body     entities.CreateAuth  true  "Login body"
// @Success      200  {object}  entities.AuthResponse
// @Failure      400  {object}  apperr.AppError
// @Failure      500  {object}  apperr.AppError
// @Router       /signin [post]
func (ac *AuthController) SignIn(c echo.Context) error {
	payload := new(entities.CreateAuth)
	if err := c.Bind(&payload); err != nil {
		return apperr.New("100002", http.StatusBadRequest, "Invalid body", "Invalid body", err)
	}

	authRes, err := ac.Authenticator.SignIn(*payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, authRes)
}
