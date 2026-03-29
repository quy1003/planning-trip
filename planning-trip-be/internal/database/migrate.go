package database

import (
	"planning-trip-be/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
		return err
	}

	return db.AutoMigrate(
		&model.User{},
		&model.Trip{},
		&model.TripMember{},
		&model.Place{},
		&model.TripPlace{},
		&model.TripScheduleItem{},
		&model.Album{},
		&model.Photo{},
		&model.PhotoTag{},
		&model.Comment{},
		&model.Reaction{},
		&model.Notification{},
	)
}
