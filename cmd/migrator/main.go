package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"just_for_fun/internal/config"
	"just_for_fun/internal/storage/db/models"
)

const (
	levelsUpToRootDir = 3
)

func main() {
	cfg := config.Load(levelsUpToRootDir)

	db, err := gorm.Open(postgres.Open(cfg.DB.ConnString()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = runMigrations(db)
	if err != nil {
		panic(err)
	}
}

func runMigrations(db *gorm.DB) error {
	_ = dropColumns(db)
	dropTables(db)

	err := db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return err
	}

	return nil
}

func dropColumns(db *gorm.DB) error {
	return nil
}

func dropTables(db *gorm.DB) {}
