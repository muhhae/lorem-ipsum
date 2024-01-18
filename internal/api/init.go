package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhhae/lorem-ipsum/internal/api/comment"
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
	e.Use(middleware.Gzip())

	e.Validator = &CustomValidator{validator: validator.New()}

	g := e.Group("/api/v1") // /api/v1

	g.Static("/static", "internal/static")
	InitUser(g)
	InitPost(g)
	InitImage(g)
	initReaction(g)
	initComment(g)
}

func InitPost(g *echo.Group) {
	p := g.Group("/post")
	p.POST("/upload", post.Upload, Auth)
	p.GET("/Default", post.Default, SoftAuth)
}

func initComment(g *echo.Group) {
	c := g.Group("/comment")
	c.POST("/send/:id", comment.SendComment, Auth)
	c.GET("/get/:id", comment.GetComment)
	c.GET("/count/:id", comment.GetCommentCount)
}

func initReaction(g *echo.Group) {
	r := g.Group("/reaction")
	r.GET("/count/:id", post.ReactionCount)
	r.GET("/myreaction/:id", post.MyReaction, SoftAuth)
	r.POST("/react/:id", post.React, Auth)
}

func InitImage(g *echo.Group) {
	i := g.Group("/image")
	i.GET("/:id", post.GetImage)
}

func InitUser(g *echo.Group) {
	u := g.Group("/user")
	u.POST("/register", user.SignUp)
	u.POST("/login", user.SignIn)
	u.GET("/session", user.Session, Auth)
	u.GET("/logout", user.SignOut)
	u.GET("/me", user.Me, Auth)
	u.GET("/myName", user.MyName, Auth)
}
