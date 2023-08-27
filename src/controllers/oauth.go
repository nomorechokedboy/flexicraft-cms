package controllers

import (
	"api/src/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	OAuthController struct {
		usecases.OAuth
	}
	OAuthQuery struct {
		Code     string `query:"code"`
		Provider string `query:"provider"`
	}
)

func NewOAuthController(usecase usecases.OAuth) *OAuthController {
	return &OAuthController{usecase}
}

func (oc *OAuthController) AuthWithOAuth2(c echo.Context) error {
	query := new(OAuthQuery)
	if err := c.Bind(query); err != nil {
		return echo.ErrBadRequest
	}

	token, err := oc.Execute(query.Code, query.Provider)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}
