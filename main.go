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

	// logger settings
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_custom} ${method} ${uri} ${status} ${error} ` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.Use(middleware.Recover())
	e.HideBanner = true

	e.Use(middleware.CORS())

	// Routes
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/login", handlers.Login)
	e.POST("/register", handlers.Register)
	e.POST("/updateuser", handlers.UpdateUser)
	e.POST("/updateuseravatar", handlers.UpdateUserAvatar)
	e.GET("/getuser", handlers.GetUser)


	e.POST("/uploadfile", handlers.UploadFile)
	e.GET("/openpvchat", handlers.OpenPvChat)

	e.GET("/ws", handlers.HandleWebSocket)

	

	log.Fatal(e.Start(":8080"))
}
