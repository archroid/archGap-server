package handlers

import (
	"archroid/archGap/db"
	"net/http"

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
