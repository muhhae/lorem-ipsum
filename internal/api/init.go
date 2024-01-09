package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhhae/lorem-ipsum/internal/api/post"
	"github.com/muhhae/lorem-ipsum/internal/api/user"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Init(e *echo.Echo) {
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	g := e.Group("/api/v1") // /api/v1

	g.Static("/static", "internal/static")
	InitPost(g)
	InitUser(g)
}

func InitPost(g *echo.Group) {
	p := g.Group("/post")
	p.POST("/upload", post.Upload, Auth)
}

func InitUser(g *echo.Group) {
	u := g.Group("/user")
	u.POST("/register", user.SignUp)
	u.POST("/login", user.SignIn)
	u.GET("/me", func(c echo.Context) error {
		id := c.Get("id")
		if id == nil {
			return c.String(401, "Login You Cunt")
		}
		return c.String(200, "Helloo"+id.(string))
	}, Auth)
}
