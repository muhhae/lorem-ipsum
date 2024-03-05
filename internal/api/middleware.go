package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/database/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtClaims struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	jwt.StandardClaims
}

// Soft Auth implement a way to identify user without kicking them out
func SoftAuth(next echo.HandlerFunc) echo.HandlerFunc {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET not set")
	}

	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			c.Set("id", primitive.NilObjectID)
			return next(c)
		}
		if cookie.Value == "" {
			c.Set("id", primitive.NilObjectID)
			return next(c)
		}

		claim := &JwtClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claim, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || token == nil || !token.Valid {
			c.Set("id", primitive.NilObjectID)
			return next(c)
		}
		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			c.Set("id", primitive.NilObjectID)
			return next(c)

		}
		u, err := user.FindOne(bson.M{"_id": claims.ID})
		if err != nil || u == nil || u.Access == user.Banned {
			c.Set("id", primitive.NilObjectID)
			return next(c)
		}
		c.Set("id", claims.ID)
		return next(c)
	}
}

// Its HARD Auth if not authenticated/authorized Get the fuck Out
func Auth(next echo.HandlerFunc) echo.HandlerFunc {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET not set")
	}

	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		if cookie.Value == "" {
			return c.NoContent(http.StatusUnauthorized)
		}

		claim := &JwtClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claim, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || token == nil || !token.Valid {
			return c.String(http.StatusUnauthorized, fmt.Sprintln("nil :", token == nil, "err : ", err.Error(), "valid :", token.Valid))
		}
		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			return c.String(http.StatusUnauthorized, "Claims not Ok")

		}
		u, err := user.FindOne(bson.M{"_id": claims.ID})
		if err != nil || u == nil || u.Access == user.Banned {
			return c.String(http.StatusUnauthorized, "User Not Ok")
		}
		c.Set("id", claims.ID)
		return next(c)
	}
}
