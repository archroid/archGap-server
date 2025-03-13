package main

import (
	"archroid/archGap/db"
	"archroid/archGap/handlers"
	"net/http"

	"github.com/charmbracelet/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db.InitDB()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// WebSocket Route
	// e.GET("/ws", echo.WrapHandler(http.HandlerFunc(websocket.HandleWebSocket)))
	// e.GET("/ws", handlers.HandleWebSocket)

	// Start WebSocket broadcasting
	// go websocket.BroadcastMessages()

	// Routes
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/login", handlers.Login)
	e.POST("/register", handlers.Register)
	e.POST("/updateuser", handlers.UpdateUser)
	e.POST("/updateuseravatar", handlers.UpdateUserAvatar)
	e.GET("/getuser", handlers.GetUser)



	e.GET("/ws", handlers.HandleWebSocket)


	log.Fatal(e.Start(":8080"))
}
