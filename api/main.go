package api

import (
	"fmt"
	"github.com/lyneq/mailapi/internal/middleware"
	"github.com/lyneq/mailapi/internal/session"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lyneq/mailapi/api/auth"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// InitForTesting initializes the Echo instance for testing purposes
func InitForTesting(e *echo.Echo) {
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1323", "http://localhost", "http://127.0.0.1:1323", "http://127.0.0.1", "https://localhost:1323", "https://localhost", "https://127.0.0.1:1323", "https://127.0.0.1"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderSetCookie, echo.HeaderCookie, "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{echo.HeaderSetCookie, "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	e.Use(session.Middleware())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Validator = &CustomValidator{validator: validator.New()}

	registerRoutes(e)
}

func Init() {
	e := echo.New()

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1323", "http://localhost", "http://127.0.0.1:1323", "http://127.0.0.1", "https://localhost:1323", "https://localhost", "https://127.0.0.1:1323", "https://127.0.0.1"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderSetCookie, echo.HeaderCookie, "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{echo.HeaderSetCookie, "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	e.Use(session.Middleware())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Validator = &CustomValidator{validator: validator.New()}

	registerRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}

// registerRoutes registers all the routes for the API
func registerRoutes(e *echo.Echo) {
	var routes []*auth.Controller
	routes = auth.GetAuthController()

	for _, route := range routes {
		if route.Active {
			switch route.Method {
			case http.MethodGet:
				if route.RequiredAuth {
					e.GET(route.Route, route.Handler, middleware.RequireAuth)
				} else {
					e.GET(route.Route, route.Handler)
				}
			case http.MethodPost:
				if route.RequiredAuth {
					e.POST(route.Route, route.Handler, middleware.RequireAuth)
				} else {
					e.POST(route.Route, route.Handler)
				}
			case http.MethodPut:
				if route.RequiredAuth {
					e.PUT(route.Route, route.Handler, middleware.RequireAuth)
				} else {
					e.PUT(route.Route, route.Handler)
				}
			case http.MethodDelete:
				if route.RequiredAuth {
					e.DELETE(route.Route, route.Handler, middleware.RequireAuth)
				} else {
					e.DELETE(route.Route, route.Handler)
				}
			default:
				fmt.Printf("Méthode non supportée ou introuvable pour %v", route.Route)
			}
		}
	}
}
