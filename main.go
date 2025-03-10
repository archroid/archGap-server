package main

import (
	"archroid/archGap/config"
	"archroid/archGap/handlers"
	"archroid/archGap/websocket"
	"net/http"

	"github.com/charmbracelet/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.InitDB()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// WebSocket Route
	e.GET("/ws", echo.WrapHandler(http.HandlerFunc(websocket.HandleWebSocket)))

	// Start WebSocket broadcasting
	go websocket.BroadcastMessages()

	// Routes
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/login", handlers.Login)
	e.POST("/register", handlers.Register)
	e.POST("/updateprofile", handlers.UpdateProfile)

	log.Fatal(e.Start(":8080"))
}
