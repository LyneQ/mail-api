package api

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/api/auth"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Init() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Validator = &CustomValidator{validator: validator.New()}

	routes := auth.GetAuthController()

	for _, route := range routes {
		if route.Active {
			switch route.Method {
			case http.MethodGet:
				e.GET(route.Route, route.Handler)
			case http.MethodPost:
				e.POST(route.Route, route.Handler)
			case http.MethodPut:
				e.PUT(route.Route, route.Handler)
			case http.MethodDelete:
				e.DELETE(route.Route, route.Handler)
			default:
				fmt.Printf("Method not supported or not found for %v", route.Route)
			}
		}
	}
	e.Logger.Fatal(e.Start(":1323"))
}
