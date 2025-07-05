package session

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var Manager *scs.SessionManager

// Init initializes the session manager with the given database connection and secure flag for cookie settings.
func Init(db *gorm.DB, secure bool) {
	Manager = scs.New()

	Manager.Lifetime = 24 * time.Hour

	Manager.Store = memstore.New()

	Manager.Cookie.Name = "session_id"
	Manager.Cookie.HttpOnly = true
	Manager.Cookie.Secure = secure
	Manager.Cookie.Path = "/"
	Manager.Cookie.SameSite = http.SameSiteLaxMode

	fmt.Println("Session manager initialized with memstore")
}

// Middleware creates an Echo middleware function that wraps requests to manage session state using the session manager.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if Manager == nil {

				fmt.Println("Warning: Session manager not initialized, using fallback initialization")
				Manager = scs.New()

				Manager.Lifetime = 24 * time.Hour
				Manager.Store = memstore.New()
				Manager.Cookie.Name = "session_id"
				Manager.Cookie.HttpOnly = true
				Manager.Cookie.Secure = false
				Manager.Cookie.Path = "/"
				Manager.Cookie.SameSite = http.SameSiteLaxMode

				fmt.Println("Session manager initialized with memstore (fallback)")
			}

			rw := c.Response().Writer

			crw := &cookieResponseWriter{
				ResponseWriter: rw,
				context:        c,
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)

				err := next(c)
				if err != nil {
					c.Error(err)
				}

				crw.processCookies()
			})

			Manager.LoadAndSave(nextHandler).ServeHTTP(crw, c.Request())

			return nil
		}
	}
}

// cookieResponseWriter is a custom response writer that captures Set-Cookie headers
// and adds them to the Echo context
type cookieResponseWriter struct {
	http.ResponseWriter
	context echo.Context
	written bool
}

// WriteHeader captures the status code and passes it to the underlying ResponseWriter
func (crw *cookieResponseWriter) WriteHeader(statusCode int) {
	if !crw.written {
		crw.processCookies()
		crw.written = true
	}
	crw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response body and passes it to the underlying ResponseWriter
func (crw *cookieResponseWriter) Write(b []byte) (int, error) {
	if !crw.written {
		crw.processCookies()
		crw.written = true
	}
	return crw.ResponseWriter.Write(b)
}

// Header captures the headers and returns them
func (crw *cookieResponseWriter) Header() http.Header {
	return crw.ResponseWriter.Header()
}

// processCookies extracts Set-Cookie headers and adds them to the Echo context
func (crw *cookieResponseWriter) processCookies() {
	h := crw.ResponseWriter.Header()

	// Check for Set-Cookie headers and add them to the Echo context
	if cookies := h["Set-Cookie"]; len(cookies) > 0 {
		for _, cookie := range cookies {
			header := http.Header{}
			header.Add("Set-Cookie", cookie)
			request := http.Request{Header: header}

			for _, parsedCookie := range request.Cookies() {
				crw.context.SetCookie(parsedCookie)
			}
		}
	}
}
