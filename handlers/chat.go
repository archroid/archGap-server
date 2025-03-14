package handlers

import (
	"archroid/archGap/db"
	"archroid/archGap/utils"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

func OpenPvChat(c echo.Context) error {
	// Get the user ID from the JWT token (you should pass the token in the Authorization header)
	userID, err := utils.ParseJWT(c.Request().Header.Get("Authorization"))

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid or expired token",
		})
	}

	var updateRequest struct {
		User2ID uint `json:"user2ID"`
	}

	if err := c.Bind(&updateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid input",
		})
	}

	if userID == updateRequest.User2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "You can't open chat with yourself",
		})

	}

	user2, _ := db.GetUser(updateRequest.User2ID)
	if user2 == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "User not found",
		})
	}

	isChatExists, chatID, _ := db.IsChatExist(userID, updateRequest.User2ID)

	if isChatExists {
		return c.JSON(http.StatusOK, map[string]uint{
			"chatID": chatID,
		})
	} else {
		chat, err := db.NewChat("pvchat", false)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal server error",
			})
		}

		_, err = db.AddUserToChat(chat.ID, []uint{userID, updateRequest.User2ID})
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal server error",
			})
		}
		return c.JSON(http.StatusOK, map[string]uint{
			"chatID": chat.ID,
		})
	}
}
