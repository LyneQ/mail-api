package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/db"
	_ "gorm.io/gorm"
	"net/http"
)

func signUpView(c echo.Context) error {

	user := new(db.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "1")
	}

	return c.JSON(http.StatusCreated, user)
}

func signInView(c echo.Context) error {
	return c.String(200, "Sign In")
}
