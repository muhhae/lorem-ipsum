package internal

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/api"
	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"github.com/muhhae/lorem-ipsum/internal/views"
)

func InitAll() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	defer connection.Disconnect(connection.Init())
	echoInit()
}

func echoInit() {
	PORT := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		PORT = ":" + p
	}
	e := echo.New()
	
	views.Init(e)
	api.Init(e)

	e.Logger.Fatal(e.Start(PORT))
}
