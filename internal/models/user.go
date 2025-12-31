package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID 			uint			`gorm:"primarykey" json:"id"`
	Name 		string			`gorm:"not null" json:"name"`
	Email 		string			`gorm:"uniqueIndex;not null" json:"email"`
	Password 	string			`gorm:"not null" json:"-"`
	CreatedAt 	time.Time		`json:"created_at"`
	UpdatedAt 	time.Time		`json:"updated_at"`
	DeletedAt 	gorm.DeletedAt	`gorm:"index" json:"-"`
}