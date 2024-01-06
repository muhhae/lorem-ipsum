package echotempl

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Templ(c echo.Context, code int, t templ.Component) error {
	c.Response().WriteHeader(code)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(c.Request().Context(), c.Response().Writer)
}
