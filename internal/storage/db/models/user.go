package models

type User struct {
	Base

	TgID      int64 `gorm:"unique"`
	UserName  *string
	FirstName *string
	LastName  *string
}
