package view

import (
	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/views"
	echotempl "github.com/muhhae/lorem-ipsum/pkg/echoTempl"
)

func Init(e *echo.Echo) {
	home(e)
	login(e)
}

func home(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
}

func login(e *echo.Echo) {
	e.GET("/login", func(c echo.Context) error {
		return echotempl.Templ(c, 200, views.LoginPage())
	})
}
