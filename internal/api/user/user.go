package user

import (
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
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	u, _ := user.FindOne(user.User{Email: req.Email})
	if u != nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "email already exists"})
	}
	u, _ = user.FindOne(user.User{Username: req.Username})
	if u != nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "username already exists"})
	}
	newUser := user.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
	_, err := newUser.Save()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	jwt, err := newUser.GenerateJWT()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
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
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	var authorized bool
	var u *user.User
	_, err := mail.ParseAddress(req.EmailOrUsername)
	if err == nil {
		u, err = user.FindOne(user.User{Username: req.EmailOrUsername})
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "user not found"})
		}
		authorized, err = u.Authenticate(req.Password)
		if !authorized || err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "wrong password"})
		}
	} else {
		u, err = user.FindOne(user.User{Email: req.EmailOrUsername})
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "user not found"})
		}
		authorized, err = u.Authenticate(req.Password)
		if !authorized || err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "wrong password"})
		}
	}
	jwt, err := u.GenerateJWT()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	c.SetCookie(&http.Cookie{
		Name:    "jwt",
		Value:   jwt,
		Expires: time.Now().Add(24 * time.Hour),
	})
	return c.NoContent(http.StatusOK)
}
