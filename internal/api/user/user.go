package user

import (
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/database/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
	count, err := user.Count(bson.M{"email": req.Email})
	fmt.Println("email", count, err)
	if count > 0 {
		return c.String(http.StatusConflict, "User with that Email already exists try sign in")
	}
	count, err = user.Count(bson.M{"username": req.Username})
	fmt.Println("username", count, err)
	if count > 0 {
		return c.String(http.StatusConflict, "Username already taken please choose another one")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	req.Password = string(hashedPassword)
	newUser := user.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
	_, err = newUser.Save()
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
	c.Response().Writer.Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusCreated)
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
	if err != nil || email == nil {
		u, err = user.FindOne(bson.M{"username": req.EmailOrUsername})
		if err != nil {
			return c.String(http.StatusNotFound, "No user found with the provided username, please check your username and try again or sign up")
		}
		authorized, err = u.Authenticate(req.Password)
		if !authorized || err != nil {
			return c.String(http.StatusUnauthorized, "Wrong password, please try again")
		}
	} else {
		u, err = user.FindOne(bson.M{"email": req.EmailOrUsername})
		if err != nil {
			return c.String(http.StatusNotFound, "No user found with the provided email, please check your email and try again or sign up")
		}
		authorized, err = u.Authenticate(req.Password)
		if !authorized || err != nil {
			return c.String(http.StatusUnauthorized, "Wrong password, please try again")
		}
	}
	jwt, err := u.GenerateJWT()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error, please try again later")
	}

	c.SetCookie(&http.Cookie{
		Name:    "jwt",
		Value:   jwt,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})
	c.Response().Writer.Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusCreated)
}

func SignOut(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "jwt",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})
	c.Response().Writer.Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}

func Session(c echo.Context) error {
	id := c.Get("id")
	if id == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	me, err := user.FindOne(bson.M{"_id": id})
	if err != nil || me == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func Me(c echo.Context) error {
	id := c.Get("id")
	if id == nil || id == "" || id == primitive.NilObjectID {
		return c.NoContent(http.StatusUnauthorized)
	}
	me, err := user.FindOne(bson.M{"_id": id})
	if err != nil || me == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	me.Password = ""
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Email":    me.Email,
		"Username": me.Username,
	})
}

func MyName(c echo.Context) error {
	id := c.Get("id")
	if id == nil || id == "" || id == primitive.NilObjectID {
		return c.NoContent(http.StatusUnauthorized)
	}
	me, err := user.FindOne(bson.M{"_id": id})
	if err != nil || me == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	me.Password = ""
	return c.String(http.StatusOK, me.Username)
}
