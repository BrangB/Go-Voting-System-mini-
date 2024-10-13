package controllers

import (
	"net/http"

	"github.com/brangb/go_voting_system/config"
	"github.com/brangb/go_voting_system/models"
	"github.com/gin-gonic/gin"
)

func GetVoteResult(c *gin.Context) {

	pollId := c.Param("id")

	var options []models.Option

	if err := config.DB.Where("poll_id = ?", pollId).Find(&options).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve voting data",
		})
		return
	}

	var voteResults []map[string]interface{}

	for _, option := range options {
		voteResults = append(voteResults, map[string]interface{}{
			"optionTitle": option.Title,
			"totalVote":   option.TotalVotes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"voteResults": voteResults,
	})

}

func VoteOption(c *gin.Context) {

	optionId := c.Param("id")

	var option models.Option

	if err := config.DB.First(&option, optionId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Option not found",
		})
		return
	}

	option.TotalVotes += 1

	if err := config.DB.Save(&option).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update vote count",
		})
		return
	}

	var newVote models.Vote

	newVote.PollID = option.PollID
	newVote.OptionID = option.ID
	newVote.Comment = c.PostForm("comment")

	if err := config.DB.Create(&newVote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create vote",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote recorded successfully",
	})
}

func GetAllVotesByOptionId(c *gin.Context) {

	optionId := c.Param("option_id")

	var votes []models.Vote

	if err := config.DB.Where("option_id = ?", optionId).Find(&votes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve votes",
		})
		return
	}

	if len(votes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No votes found for this option",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"votes": votes,
	})
}
