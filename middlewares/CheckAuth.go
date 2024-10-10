package middlewares

import (
	"net/http"

	"github.com/brangb/go_voting_system/config"
	"github.com/brangb/go_voting_system/models"
	"github.com/brangb/go_voting_system/utils"
	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {

	tokenString, err := c.Cookie("Access_Token")
	if err != nil {
		refreshTokenString, err := c.Cookie("Refresh_Token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized. No valid tokens provided.",
			})
			c.Abort()
			return
		}

		_, claims, err := utils.ValidateToken(refreshTokenString, "refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized. Invalid or expired refresh token.",
			})
			c.Abort()
			return
		}

		var user models.User
		config.DB.Preload("Poll").First(&user, claims["User_ID"])
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized. User not found.",
			})
			c.Abort()
			return
		}

		newAccessToken, err := utils.GenerateAccessToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not generate a new access token.",
			})
			c.Abort()
			return
		}

		c.SetCookie("Access_Token", newAccessToken, 15*60, "", "", false, true)
		c.Set("user", user)
		c.Next()
		return
	}

	_, claims, err := utils.ValidateToken(tokenString, "access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	var user models.User
	config.DB.Preload("Poll").First(&user, claims["User_ID"])
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized. User not found.",
		})
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
