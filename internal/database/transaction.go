package database

import (
	"gorm.io/gorm"
)

// TransactionManager - interface untuk transaction operator
type TransactionManager interface {
	WithTransaction(fn func(tx *gorm.DB) error) error
}

// transactionManager - implementasi transaction manager
type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) WithTransaction(fn func(tx *gorm.DB) error) error {
	tx := tm.db.Begin()
	defer func(){
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

