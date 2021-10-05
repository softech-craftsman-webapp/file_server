package main

import (
	"fmt"
	"os"

	bootstrap "file_server/bootstrap"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// @title File Server
// @version 1.0
// @description File API Service.

// @host 127.0.0.1:5000
// @BasePath /

// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(".env file is not imported, in production kindly ignore this message")
	}

	/*
	   |--------------------------------------------------------------------------
	   | Start Server
	   |--------------------------------------------------------------------------
	*/
	app := echo.New()
	port := os.Getenv("PORT")

	// Application
	bootstrap.Start(app, port)
}
