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

func Auth(next echo.HandlerFunc) echo.HandlerFunc {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET not set")
	}

	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			fmt.Println("1", err)
			return c.NoContent(http.StatusUnauthorized)
		}
		if cookie.Value == "" {
			fmt.Println("2", "empty cookie")
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
			fmt.Println("3", err, token)
			return c.NoContent(http.StatusUnauthorized)
		}
		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			fmt.Println("4", "claims not ok")
			return c.NoContent(http.StatusUnauthorized)

		}
		u, err := user.FindOne(bson.M{"_id": claims.ID})
		if err != nil || u == nil || u.Access == user.Banned {
			fmt.Println("5", "User not found or banned")
			return c.NoContent(http.StatusUnauthorized)
		}
		fmt.Println("6", "Success")
		c.Set("id", claims.ID)
		return next(c)
	}
}
