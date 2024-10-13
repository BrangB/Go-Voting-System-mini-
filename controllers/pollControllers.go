package controllers

import (
	"net/http"
	"time"

	"github.com/brangb/go_voting_system/config"
	"github.com/brangb/go_voting_system/models"
	"github.com/brangb/go_voting_system/utils"
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

	// Validate & Get user data
	user, valid := utils.UserUtils(c)

	if !valid {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
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
		OwnerID:     user.ID,
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

func GetAllPolls(c *gin.Context) {

	user, valid := utils.UserUtils(c)

	if !valid {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var polls []models.Poll
	result := config.DB.Where("owner_id = ?", user.ID).Preload("Options").Find(&polls)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve polls",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"polls": polls,
	})

}

func GetPollById(c *gin.Context) {

	pollID := c.Param("id")

	var poll models.Poll

	config.DB.Preload("Options").First(&poll, pollID)

	if poll.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "There is no poll",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"poll": poll,
	})

}

func DeletePollByID(c *gin.Context) {

	pollId := c.Param("id")

	if pollId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing Poll ID",
		})
		return
	}

	var deletePoll models.Poll

	if err := config.DB.First(&deletePoll, pollId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Poll not found",
		})
		return
	}

	if err := config.DB.Delete(&deletePoll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete the poll",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poll deleted successfully",
		"poll":    deletePoll,
	})
}

func UpdatePollByID(c *gin.Context) {
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

	var pollBody PollBody

	if err := c.ShouldBindJSON(&pollBody); err != nil {
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

	// Get poll ID
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing Poll ID",
		})
		return
	}

	var poll models.Poll
	if err := config.DB.Preload("Options").First(&poll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Poll not found",
		})
		return
	}

	if poll.OwnerID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to update this poll",
		})
		return
	}

	poll.Title = pollBody.Title
	poll.Description = pollBody.Description
	poll.ImgUrl = pollBody.ImgUrl

	startDate, err := time.Parse("2006-01-02", pollBody.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start date format",
		})
		return
	}
	endDate, err := time.Parse("2006-01-02", pollBody.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end date format",
		})
		return
	}
	poll.StartDate = startDate
	poll.EndDate = endDate

	poll.Status = pollBody.Status
	poll.Public = pollBody.Public

	config.DB.Where("poll_id = ?", poll.ID).Delete(&models.Option{})

	var newOptions []models.Option
	for _, updatedOption := range pollBody.Options {
		newOption := models.Option{
			Title:  updatedOption.Title,
			ImgUrl: updatedOption.ImgUrl,
			PollID: poll.ID,
		}
		newOptions = append(newOptions, newOption)
	}

	poll.Options = newOptions

	if err := config.DB.Save(&poll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update poll",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poll and options updated successfully",
		"poll":    poll,
	})
}
