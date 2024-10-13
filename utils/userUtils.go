package utils

import (
	"net/http"

	"github.com/brangb/go_voting_system/models"
	"github.com/gin-gonic/gin"
)

func UserUtils(c *gin.Context) (models.User, bool) {

	userData, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return models.User{}, false
	}

	user, ok := userData.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to retrieve user data",
		})
		return models.User{}, false
	}

	return user, true

}
