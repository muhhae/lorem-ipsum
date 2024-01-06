package router

import (
	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/router/view"
)

func Init(e *echo.Echo) {
	view.Init(e)
	e.Static("/static", "internal/static")
}
