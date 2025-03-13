package handlers

import (
	"archroid/archGap/config"
	"archroid/archGap/models"
	"archroid/archGap/services"

	"github.com/labstack/echo/v4"

	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Handle file upload
func HandleFileUpload(c echo.Context) error {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "No file uploaded")
	}

	// Save file to the local storage (e.g., "uploads/")
	dst := filepath.Join("uploads", file.Filename)
	out, err := os.Create(dst)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file")
	}
	defer out.Close()

	// Copy the uploaded file to the destination
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to copy file")
	}


	err = services.SaveMessage(config.DB, models.Message{
		SenderID:   1, // Replace with the actual user ID
		Content:    file.Filename,
		MessageType: "file",

	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save attachment")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "File uploaded successfully",
		"file":    file.Filename,
	})
}
