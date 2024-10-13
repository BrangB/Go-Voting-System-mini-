package controllers

import (
	"net/http"

	"github.com/brangb/go_voting_system/config"
	"github.com/brangb/go_voting_system/models"
	"github.com/brangb/go_voting_system/utils"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func RegisterUser(c *gin.Context) {
	var Body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBind(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to read data",
		})
		return
	}

	if Body.Username == "" || Body.Email == "" || Body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The body can't be empty.",
		})
	}

	hashedPassword, err := utils.HashPassword(Body.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to hash your password",
		})
	}

	user := models.User{
		Username: Body.Username,
		Email:    Body.Email,
		Password: hashedPassword,
	}

	result := config.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to create your account",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account is created successfully!!",
	})

}

func Login(c *gin.Context) {
	var Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBind(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to read data the body",
		})
		return
	}

	if Body.Email == "" || Body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The body can't be empty",
		})
	}

	var user models.User

	config.DB.First(&user, "email = ?", Body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid User",
		})
		return
	}

	IsRightPassword := utils.CheckHashedPassword(Body.Password, user.Password)

	if !IsRightPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong Password",
		})
		return
	}

	access_token, err := utils.GenerateAccessToken(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to generate Token.",
		})
	}

	refresh_token, err := utils.GenerateRefreshToken(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to generate Token.",
		})
	}

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("Access_Token", access_token, 15*60, "", "", false, true)
	c.SetCookie("Refresh_Token", refresh_token, 30*24*60*60, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})

}

func Logout(c *gin.Context) {

	// Remove Access Token, Refresh Token and user data
	c.SetCookie("Access_Token", "", -1, "/", "", false, true)
	c.SetCookie("Refresh_Token", "", -1, "/", "", false, true)
	c.Set("user", nil)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func GetUserProfile(c *gin.Context) {

	// Validate & Get user data
	user, valid := utils.UserUtils(c)

	if !valid {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var userProfile models.User

	err := config.DB.Preload("Poll").First(&userProfile, user.ID).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userData": userProfile,
	})

}

func UpdateProfile(c *gin.Context) { //update username or email

	var Body struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body or missing required fields",
		})
		return
	}

	// Validate & Get user data
	user, valid := utils.UserUtils(c)

	if !valid {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var userProfile models.User

	if err := config.DB.First(&userProfile, user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userProfile.Username = Body.Username
	userProfile.Email = Body.Email

	if err := config.DB.Save(&userProfile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fail to update user info."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User's info is updated successfully.",
	})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
