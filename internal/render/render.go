package render

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Render renders a templ component to the response
func Render(c echo.Context, status int, t templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(status)
	return t.Render(c.Request().Context(), c.Response().Writer)
}
