package config

import (
	"fmt"

	"github.com/brangb/go_voting_system/models"
)

func SyncDatabases() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Poll{},
		&models.Option{},
		&models.Vote{},
	)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database migration is successfully done!!!")
	}
}
