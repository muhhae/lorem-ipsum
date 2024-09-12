package internal

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/api"
	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"github.com/muhhae/lorem-ipsum/internal/views"
)

func InitAll() {
	defer connection.Disconnect(connection.Init())
	echoInit()
}

func echoInit() {
	PORT := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		PORT = ":" + p
	}
	e := echo.New()
	defer e.Close()

	views.Init(e)
	api.Init(e)
	e.Logger.Fatal(e.Start(PORT))
}
