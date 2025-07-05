package middleware // protectedRoute handles requests to protected endpoints by checking the user's authentication status via a session.

import (
	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/internal/session"
	"net/http"
)

// It returns unauthorized status if the user is not logged in, otherwise responds with the user's data.
func protectedRoute(c echo.Context) error {
	user, err := session.GetCurrentUser(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Not logged in",
		})
	}
	return c.JSON(200, user)
}

// RequireAuth is a middleware that enforces authentication by verifying a valid user session before proceeding.
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := session.GetUserID(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Authentication required",
			})
		}
		return next(c)
	}
}
