package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	Route   string
	Method  string
	Active  bool
	Handler func(c echo.Context) error
}

func GetAuthController() []*Controller {
	return []*Controller{
		{
			Route:   "/signup",
			Method:  http.MethodPost,
			Active:  true,
			Handler: signUpView,
		},
		{
			Route:   "/signin",
			Method:  http.MethodPost,
			Active:  true,
			Handler: signInView,
		},
	}
}
