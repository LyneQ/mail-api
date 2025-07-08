package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/config"
	"github.com/lyneq/mailapi/db"
	"github.com/lyneq/mailapi/internal/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

// isAllowedDomain checks if the given URL's domain is in the list of allowed domains
func isAllowedDomain(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		return false
	}

	host := parsedURL.Hostname()

	allowedDomains := config.GetAllowedDomains()
	for _, domain := range allowedDomains {
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
}

type SignUpRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=64"`
	Password    string `json:"password" validate:"required,min=3,max=64"`
	CallbackURL string `json:"callbackURL"`
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

	// Check for callbackURL in query parameters (this takes precedence over JSON body)
	queryCallbackURL := c.QueryParam("callbackURL")
	fmt.Printf("Query callbackURL (signup): %v\n", queryCallbackURL)
	if queryCallbackURL != "" {
		req.CallbackURL = queryCallbackURL
		fmt.Printf("Using callbackURL from query parameter (signup): %v\n", req.CallbackURL)
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

	// Check if there's a callback URL and if it's allowed
	if req.CallbackURL != "" {
		if isAllowedDomain(req.CallbackURL) {
			return c.Redirect(http.StatusSeeOther, req.CallbackURL)
		} else {
			// If the domain is not allowed, log it and continue without redirection
			_ = fmt.Errorf("redirection to unauthorized domain attempted: %s", req.CallbackURL)
		}
	}

	return c.JSON(http.StatusCreated, user)
}

type SignInRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=64"`
	Password    string `json:"password" validate:"required,min=3,max=64"`
	CallbackURL string `json:"callbackURL"`
}

// signInView handles user sign-in by validating credentials and creating a session.
func signInView(c echo.Context) error {
	var req SignInRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "An error occurred while processing your request.",
		})
	}

	// Check for callbackURL in query parameters (this takes precedence over JSON body)
	queryCallbackURL := c.QueryParam("callbackURL")
	if queryCallbackURL != "" {
		req.CallbackURL = queryCallbackURL
		fmt.Printf("Using callbackURL from query parameter: %v\n", req.CallbackURL)
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

	// Check if there's a callback URL and if it's allowed
	if req.CallbackURL != "" {
		fmt.Printf("CallbackURL is not empty, checking if domain is allowed\n")
		if isAllowedDomain(req.CallbackURL) {
			fmt.Printf("Domain is allowed, redirecting to: %v\n", req.CallbackURL)
			return c.Redirect(http.StatusSeeOther, req.CallbackURL)
		} else {
			// If the domain is not allowed, log it and continue without redirection
			_ = fmt.Errorf("redirection to unauthorized domain attempted: %s", req.CallbackURL)
		}
	}

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

// signOutView handles user sign-out by invalidating the session and deleting the session cookie.
func signOutView(c echo.Context) error {
	// Get callbackURL from query parameter
	callbackURL := c.QueryParam("callbackURL")
	fmt.Printf("Query callbackURL (signout): %v\n", callbackURL)

	err := session.Manager.Destroy(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "An error occurred while signing out.",
			"error":   err.Error(),
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = session.Manager.Cookie.Name
	cookie.Value = ""
	cookie.Path = session.Manager.Cookie.Path
	cookie.Expires = time.Now().Add(-1 * time.Hour) // Set expiry in the past
	cookie.MaxAge = -1                              // Delete the cookie
	cookie.Secure = session.Manager.Cookie.Secure
	cookie.HttpOnly = session.Manager.Cookie.HttpOnly
	cookie.SameSite = session.Manager.Cookie.SameSite

	c.SetCookie(cookie)

	// Check if there's a callback URL and if it's allowed
	if callbackURL != "" {
		if isAllowedDomain(callbackURL) {
			return c.Redirect(http.StatusSeeOther, callbackURL)
		} else {
			// If the domain is not allowed, log it and continue without redirection
			_ = fmt.Errorf("redirection to unauthorized domain attempted: %s", callbackURL)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully signed out.",
	})
}
