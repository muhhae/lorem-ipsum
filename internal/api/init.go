package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhhae/lorem-ipsum/internal/api/post"
)

func Init(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("/api/v1")

	g.Static("/static", "internal/static")

}

func InitPost(g *echo.Group) {
	p := g.Group("/post")
	p.POST("/upload", post.Upload, Auth)
}
