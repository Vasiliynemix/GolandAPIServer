package models

type Base struct {
	ID uint `gorm:"primaryKey; autoIncrement"`

	createdAt int64 `gorm:"autoCreateTime:milli"`
	updatedAt int64 `gorm:"autoUpdateTime:milli"`
}
