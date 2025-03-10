package handlers

import (
	"archroid/archGap/config"
	"archroid/archGap/services"
	"archroid/archGap/utils"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	// Get email and password from the request body
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid input",
		})
	}

	// Call the login service
	user, token, err := services.LoginUser(config.DB, loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	// Return the JWT token and user info (optional)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
		"token": token,
	})
}

func Register(c echo.Context) error {
	// Get email and password from the request body
	var registerRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind request to struct
	if err := c.Bind(&registerRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid input",
		})
	}

	println("eml"+registerRequest.Email, "pss"+registerRequest.Password)

	// Call the registration service
	user, err := services.RegisterUser(config.DB, registerRequest.Email, registerRequest.Password)
	if err != nil {
		// Handle errors: this might include email already in use, etc.
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	// Return the created user (excluding password) as a response
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func UpdateProfile(c echo.Context) error {
	// Get the user ID from the JWT token (you should pass the token in the Authorization header)
	userID, err := utils.ParseJWT(c.Request().Header.Get("Authorization"))

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid or expired token",
		})
	}

	// Get updated profile data from the request body
	var updateRequest struct {
		Username       string `json:"username"`
		ProfilePicture string `json:"profilepicture"`
	}

	if err := c.Bind(&updateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid input",
		})
	}

	// Call the update profile service
	user, err := services.UpdateProfile(config.DB, userID, updateRequest.Username, updateRequest.ProfilePicture)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	// Return updated user profile (without password)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func UpdateAvatar(c echo.Context) error {
	// Get the user ID from the JWT token (you should pass the token in the Authorization header)
	userID, err := utils.ParseJWT(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid or expired token",
		})
	}

	// Retrieve the file from the form data
	file, err := c.FormFile("profilepicture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Failed to get file",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to open file",
		})
	}
	defer src.Close()

	// Save the file to a specific location
	dst, err := utils.SaveFile(src, fmt.Sprintf("%d_%s", userID, filepath.Ext(file.Filename)), "profilepicture")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to save file",
		})
	}

	var username string

	// Update the user's profile picture in the database
	user, err := services.UpdateProfile(config.DB, userID, username, dst)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	// Return updated user profile (without password)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":             user.ID,
			"email":          user.Email,
			"profilepicture": user.ProfilePicture,
		},
	})
}
