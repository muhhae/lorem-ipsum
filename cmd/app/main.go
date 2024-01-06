package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/router"
)

func main() {
	PORT := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		PORT = ":" + p
	}
	e := echo.New()
	router.Init(e)

	e.Logger.Fatal(e.Start(PORT))
}
