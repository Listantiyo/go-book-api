package models

import (
	"time"

	"gorm.io/gorm"
)

type BorrowStatus string

const (
	BorrowStatusBorrowed BorrowStatus = "borrowed"
	BorrowStatusReturned BorrowStatus = "returned"
	BorrowStatusOverdue BorrowStatus = "overdue"
)

type Borrow struct {
	ID uint `gorm:"primarykey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	BookID uint `gorm:"not null;index" json:"book_id"`
	BorrowDate time.Time `gorm:"not null" json:"borrow_date"`
	DueDate time.Time `gorm:"not null" json:"due_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
	Status BorrowStatus `gorm:"type:varchar(20);check:status IN ('borrowed','returned','overdue');not null" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`


	// Relations
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Book Book `gorm:"foreignKey:BookID;references:ID" json:"book,omitempty"`
}