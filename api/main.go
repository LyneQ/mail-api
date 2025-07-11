package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lyneq/mailapi/api/auth"
	"github.com/lyneq/mailapi/api/email"
	"github.com/lyneq/mailapi/config"
	"github.com/lyneq/mailapi/internal/middleware"
	"github.com/lyneq/mailapi/internal/session"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Init() {

	allowedDomains := config.GetAllowedDomains()
	e := echo.New()

	var allowedHosts []string

	for _, domain := range allowedDomains {
		allowedHosts = append(allowedHosts, "http://"+domain)
	}

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     allowedHosts,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderSetCookie, echo.HeaderCookie, "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{echo.HeaderSetCookie, "Set-Cookie", "X-SvelteKit-Action"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	e.Use(session.Middleware())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Validator = &CustomValidator{validator: validator.New()}

	registerRoutes(e)

	// Start the main server on HTTP only
	port := config.GetAPIPort()
	if port == "" {
		port = "1323"
	}
	e.Logger.Fatal(e.Start(":" + port))

}

// registerRoutes registers all the routes for the API
func registerRoutes(e *echo.Echo) {
	var routes []*auth.Controller
	routes = auth.GetAuthController()

	emailRoutes := email.GetEmailController()
	for _, route := range emailRoutes {
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
