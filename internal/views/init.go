package views

import (
	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/views/home"
	"github.com/muhhae/lorem-ipsum/internal/views/login"
	"github.com/muhhae/lorem-ipsum/internal/views/register"
	echotempl "github.com/muhhae/lorem-ipsum/pkg/echoTempl"
)

func Init(e *echo.Echo) {
	homePage(e)
	loginPage(e)
	registerPage(e)
}

func homePage(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return echotempl.Templ(c, 200, home.Index())
	})
}

func loginPage(e *echo.Echo) {
	e.GET("/login", func(c echo.Context) error {
		return echotempl.Templ(c, 200, login.Index())
	})
}

func registerPage(e *echo.Echo) {
	e.GET("/register", func(c echo.Context) error {
		return echotempl.Templ(c, 200, register.Index())
	})
}
