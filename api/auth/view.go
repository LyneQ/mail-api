package auth

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"net/http"
)

type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=3,max=64"`
}

// signUpView handles user sign-up by validating request data, hashing passwords, and creating new user records in the database.
func signUpView(c echo.Context) error {

	// Bind the request body to the SignUpRequest struct
	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "An error occurred while processing your request.",
		})
	}

	// Validate request data
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Could not validate your information.",
		})
	}

	// Check if the username already exists
	var existingUser db.User
	result := db.DB.Where("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Something went wrong, please try again later.",
		})
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		_ = fmt.Errorf("database error: %v", result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "An unexpected error occurred, please try again later.",
		})
	}

	// Hash the password before store it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		_ = fmt.Errorf("bcrypt error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "An internal error occurred, please try again later.",
		})
	}

	// Create a new user on the database
	user := &db.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}
	if result := db.DB.Create(user); result.Error != nil {
		_ = fmt.Errorf("database error: %v", result.Error)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Something went wrong, please try again later.",
		})
	}

	return c.JSON(http.StatusCreated, user)
}

func signInView(c echo.Context) error {
	return c.String(200, "Sign In")
}
