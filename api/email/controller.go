package email

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Route        string
	Method       string
	Active       bool
	Handler      func(c echo.Context) error
	RequiredAuth bool
}

func GetEmailController() []*Controller {
	return []*Controller{
		{
			Route:        "/api/email/inbox",
			Method:       http.MethodGet,
			Active:       true,
			Handler:      getInboxView,
			RequiredAuth: true,
		},
		{
			Route:        "/api/email/folder",
			Method:       http.MethodGet,
			Active:       true,
			Handler:      getFolderView,
			RequiredAuth: true,
		},
		{
			Route:        "/api/email/:id",
			Method:       http.MethodGet,
			Active:       true,
			Handler:      getEmailView,
			RequiredAuth: true,
		},
		{
			Route:        "/api/email/send",
			Method:       http.MethodPost,
			Active:       true,
			Handler:      sendEmailView,
			RequiredAuth: true,
		},
		{
			Route:        "/api/email",
			Method:       http.MethodGet,
			Active:       true,
			Handler:      getFoldersView,
			RequiredAuth: true,
		},
	}
}
