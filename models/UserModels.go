package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique; not null"`
	Email    string `gorm:"unique; not null"`
	Password string `gorm:"not null"`
	Poll     []Poll `gorm:"foreignKey:OwnerID"`
}

type Poll struct {
	gorm.Model
	OwnerID     uint      `gorm:"not null"`
	Owner       User      `gorm:"foreignKey:OwnerID"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	ImgUrl      string    `gorm:"not null"`
	StartDate   time.Time `gorm:"not null"`
	EndDate     time.Time `gorm:"not null"`
	Status      bool      `gorm:"not null"`
	Public      bool      `gorm:"not null"`
	Options     []Option  `gorm:"constraint:OnDelete:CASCADE"`
}

type Option struct {
	gorm.Model
	PollID     uint   `gorm:"not null"`
	Poll       Poll   `gorm:"foreignKey:PollID"`
	Title      string `gorm:"not null"`
	ImgUrl     string
	TotalVotes int    `gorm:"default:0"`
	Votes      []Vote `gorm:"constraint:OnDelete:CASCADE"`
}

type Vote struct {
	gorm.Model
	PollID   uint   `gorm:"not null"`
	OptionID uint   `gorm:"not null"`
	Poll     Poll   `gorm:"foreignKey:PollID"`
	Option   Option `gorm:"foreignKey:OptionID"`
	Comment  string
}
