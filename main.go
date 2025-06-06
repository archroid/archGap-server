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

	e.POST("/api/login", handlers.Login)
	e.POST("/api/register", handlers.Register)
	e.POST("/api/updateuser", handlers.UpdateUser)
	e.POST("/api/updateuseravatar", handlers.UpdateUserAvatar)
	e.POST("/api/verifytoken", handlers.Verifytoken)
	e.GET("/api/getuser", handlers.GetUser)

	e.GET("/api/getchatsbyuser", handlers.GetChatsbyUser)

	e.POST("/api/uploadfile", handlers.UploadFile)
	e.GET("/api/openpvchat", handlers.OpenPvChat)

	e.GET("/api/ws", handlers.HandleWebSocket)

	e.Static("/static", "web") // Serve static files under /static

	e.GET("/", func(c echo.Context) error {
		return c.File("web/index.html")
	})

	e.GET("/login", func(c echo.Context) error {
		return c.File("web/login.html")
	})

	e.GET("/chat/:id", func(c echo.Context) error {
		return c.File("web/chat.html")
	})

	e.GET("/chat", func(c echo.Context) error {
		return c.File("web/chat.html")
	})

	log.Fatal(e.Start(":8080"))
}
