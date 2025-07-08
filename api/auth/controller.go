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
			Route:        "/api/signup",
			Method:       http.MethodPost,
			Active:       true,
			Handler:      signUpView,
			RequiredAuth: false,
		},
		{
			Route:        "/api/signin",
			Method:       http.MethodPost,
			Active:       true,
			Handler:      signInView,
			RequiredAuth: false,
		}, {
			Route:        "/api/me",
			Method:       http.MethodGet,
			Active:       true,
			Handler:      me,
			RequiredAuth: true,
		}, {
			Route:        "/api/signout",
			Method:       http.MethodGet,
			Active:       true,
			Handler:      signOutView,
			RequiredAuth: true,
		},
	}
}
