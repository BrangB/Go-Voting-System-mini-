package config

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("VotingSystem.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
