package controllers

import (
	"net/http"
	"time"

	"github.com/brangb/go_voting_system/config"
	"github.com/brangb/go_voting_system/models"
	"github.com/gin-gonic/gin"
)

func CreatePoll(c *gin.Context) {

	type OptionBody struct {
		Title  string `json:"title"`
		ImgUrl string `json:"img_url"`
	}

	type PollBody struct {
		Title       string       `json:"title"`
		Description string       `json:"description"`
		ImgUrl      string       `json:"img_url"`
		StartDate   string       `json:"start_date"`
		EndDate     string       `json:"end_date"`
		Status      bool         `json:"status"`
		Public      bool         `json:"public"`
		Options     []OptionBody `json:"options"`
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var pollBody PollBody

	if err := c.ShouldBindJSON(&pollBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	startDate, err := time.Parse(time.RFC3339, pollBody.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start date format",
		})
		return
	}

	endDate, err := time.Parse(time.RFC3339, pollBody.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end date format",
		})
		return
	}

	poll := models.Poll{
		OwnerID:     user.(models.User).ID,
		Title:       pollBody.Title,
		Description: pollBody.Description,
		ImgUrl:      pollBody.ImgUrl,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      pollBody.Status,
		Public:      pollBody.Public,
	}

	for _, optReq := range pollBody.Options {
		option := models.Option{
			Title:  optReq.Title,
			ImgUrl: optReq.ImgUrl,
		}
		poll.Options = append(poll.Options, option)
	}

	result := config.DB.Create(&poll)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to create new voting room.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": poll,
	})
}
