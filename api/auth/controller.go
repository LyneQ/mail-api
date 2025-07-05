package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	Route        string
	Method       string
	Active       bool
	Handler      func(c echo.Context) error
	RequiredAuth bool
}

func GetAuthController() []*Controller {
	return []*Controller{
		{
			Route:        "/signup",
			Method:       http.MethodPost,
			Active:       true,
			Handler:      signUpView,
			RequiredAuth: false,
		},
		{
			Route:        "/signin",
			Method:       http.MethodPost,
			Active:       true,
			Handler:      signInView,
			RequiredAuth: false,
		}, {
			Route:        "/me",
			Method:       http.MethodGet,
			Active:       true,
			Handler:      me,
			RequiredAuth: true,
		},
	}
}
