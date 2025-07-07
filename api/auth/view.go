package auth

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/db"
	"github.com/lyneq/mailapi/internal/session"
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
		_ = fmt.Errorf("username already exists %v\n", existingUser)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Username unavailable. Please try another one.",
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

type SignInRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=3,max=64"`
}

// signInView handles user sign-in by validating credentials and creating a session.
func signInView(c echo.Context) error {
	var req SignInRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "An error occurred while processing your request.",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Could not validate your information.",
		})
	}

	var user db.User
	result := db.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			_ = fmt.Errorf("user not found: %v", req.Username)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid username or password.",
			})
		}
		_ = fmt.Errorf("database error: %v", result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "An unexpected error occurred, please try again later.",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		_ = fmt.Errorf("invalid password for user: %v", req.Username)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid username or password.",
		})
	}

	session.SetSessionCookie(c, user.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"message":  "Login successful. Session cookie has been set.",
	})
}

func me(c echo.Context) error {
	user, err := session.GetCurrentUser(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Not authenticated. Please sign in first.",
			"error":   err.Error(),
		})
	}

	cookies := c.Cookies()
	cookieNames := make([]string, len(cookies))
	for i, cookie := range cookies {
		cookieNames[i] = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "authenticated!",
		"user":    user,
	})
}
