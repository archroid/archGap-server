package handlers

import (
	"archroid/archGap/db"
	"archroid/archGap/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

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
	// Call the registration service
	user, err := db.RegisterUser(registerRequest.Email, registerRequest.Password)
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
	user, token, err := db.LoginUser(loginRequest.Email, loginRequest.Password)
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

func UpdateUser(c echo.Context) error {
	// Get the user ID from the JWT token (you should pass the token in the Authorization header)
	userID, err := utils.ParseJWT(c.Request().Header.Get("Authorization"))

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid or expired token",
		})
	}

	// Get updated profile data from the request body
	var updateRequest struct {
		Name           string    `json:"name"`
		ProfilePicture string    `json:"profilepicture"`
		LastSeen       time.Time `json:"lastseen"`
		IsOnline       bool      `json:"isonline"`
	}

	if err := c.Bind(&updateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid input",
		})
	}

	// Call the update profile service
	err = db.UpdateUser(userID, &updateRequest.Name, &updateRequest.ProfilePicture, &updateRequest.LastSeen, &updateRequest.IsOnline)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	// Return updated user profile (without password)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "updated",
	})
}

func GetUser(c echo.Context) error {
	var input struct {
		UserID uint `json:"userID"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid input",
		})
	}
	// Call the get user service
	user, err := db.GetUser(input.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	// Return the user profile (without password)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":             user.ID,
			"email":          user.Email,
			"name":           user.Name,
			"profilepicture": user.ProfilePicture,
			"lastseen":       user.LastSeen,
			"isonline":       user.IsOnline,
		},
	})
}

func UpdateUserAvatar(c echo.Context) error {
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
	err = db.UpdateUser(userID, &username, &dst, nil, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "updated",
	})

}
