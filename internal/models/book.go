package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID 			uint			`gorm:"primarykey" json:"id"`
	Title		string			`gorm:"not null" json:"title"`
	Author		string			`gorm:"not null" json:"author"`
	ISBN		string			`gorm:"uniqueIndex" json:"isbn"`
	Description string			`gorm:"type:text" json:"description"`
	Stock		int				`gorm:"type:integer;default:0" json:"stock"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
	DeletedAt 	gorm.DeletedAt	`gorm:"index" json:"-"`
}