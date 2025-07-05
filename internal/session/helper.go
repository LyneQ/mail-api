package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/db"
	"net/http"
	"time"
)

// GetUserID extracts the user ID from the provided context and returns it. Returns an error if the user is unauthenticated.
func GetUserID(ctx context.Context) (uint, error) {
	userID, ok := Manager.Get(ctx, "userID").(uint)
	if !ok || userID == 0 {
		return 0, errors.New("unauthenticated")
	}
	return userID, nil
}

// GetCurrentUser retrieves the current authenticated user based on the user ID stored in the context.
// Returns a pointer to a User and an error if the user cannot be found or is not authenticated.
func GetCurrentUser(ctx context.Context) (*db.User, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user db.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// SetSessionCookie explicitly sets a session cookie for the given Echo context
// This is a workaround for cases where the session middleware doesn't properly set the cookie
func SetSessionCookie(c echo.Context, userID uint) {
	fmt.Println("Setting session cookie for user:", userID)

	Manager.Put(c.Request().Context(), "userID", userID)

	token, expiry, err := Manager.Commit(c.Request().Context())
	if err != nil {
		fmt.Printf("Error committing session: %v\n", err)
		return
	}

	fmt.Printf("Session token after commit: %v\n", token)
	if token == "" {
		fmt.Println("Token is still empty after commit")
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = Manager.Cookie.Name
	cookie.Value = token
	cookie.Path = Manager.Cookie.Path
	//cookie.Domain = Manager.Cookie.Domain
	cookie.Expires = expiry
	cookie.MaxAge = int(expiry.Sub(time.Now()).Seconds())
	cookie.Secure = Manager.Cookie.Secure
	cookie.HttpOnly = Manager.Cookie.HttpOnly
	cookie.SameSite = Manager.Cookie.SameSite

	c.SetCookie(cookie)

	c.Response().Header().Add("Set-Cookie", cookie.String())
}
