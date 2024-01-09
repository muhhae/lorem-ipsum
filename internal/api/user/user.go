package user

import (
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/database/user"
)

type SignUpRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

func SignUp(c echo.Context) error {
	req := SignUpRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Please fill in valid input")
	}
	if err := c.Validate(&req); err != nil {
		return c.String(http.StatusBadRequest, "Please fill in valid input")
	}
	u, _ := user.FindOne(user.User{Email: req.Email})
	if u != nil {
		return c.String(http.StatusConflict, "User with that Email already exists try sign in")
	}
	u, _ = user.FindOne(user.User{Username: req.Username})
	if u != nil {
		return c.String(http.StatusConflict, "Username already taken please choose another one")
	}
	newUser := user.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
	_, err := newUser.Save()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	jwt, err := newUser.GenerateJWT()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	c.SetCookie(&http.Cookie{
		Name:    "jwt",
		Value:   jwt,
		Expires: time.Now().Add(24 * time.Hour),
	})
	return c.NoContent(http.StatusOK)
}

type SignInRequest struct {
	EmailOrUsername string `form:"email" json:"email" validate:"required"`
	Password        string `form:"password" json:"password" validate:"required"`
}

func SignIn(c echo.Context) error {
	req := SignInRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	var authorized bool
	var u *user.User
	email, err := mail.ParseAddress(req.EmailOrUsername)
	fmt.Println(email)
	if err != nil || email == nil {
		u, err = user.FindOne(user.User{Username: req.EmailOrUsername})
		if err != nil {
			return c.String(http.StatusNotFound, "User with that Username not found")
		}
		authorized, err = u.Authenticate(req.Password)
		if !authorized || err != nil {
			return c.String(http.StatusUnauthorized, "Wrong password")
		}
	} else {
		u, err = user.FindOne(user.User{Email: req.EmailOrUsername})
		if err != nil {
			return c.String(http.StatusNotFound, "User with that Email not found")
		}
		authorized, err = u.Authenticate(req.Password)
		if !authorized || err != nil {
			return c.String(http.StatusUnauthorized, "Wrong password")
		}
	}
	jwt, err := u.GenerateJWT()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	c.SetCookie(&http.Cookie{
		Name:    "jwt",
		Value:   jwt,
		Expires: time.Now().Add(24 * time.Hour),
	})
	return c.Redirect(http.StatusFound, "/")
}
