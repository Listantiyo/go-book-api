package repository

import (
	"book-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BorrowRepository interface {
	Create(borrow *models.Borrow) error
	CreateWithTx(tx *gorm.DB, borrow *models.Borrow) error
	FindByID(id uint) (*models.Borrow, error)
	FindByIDWithLock(tx *gorm.DB ,id uint) (*models.Borrow, error)
	FindByUserID(userID uint, limit, offset int) ([]models.Borrow, error)
	Update(borrow *models.Borrow) error
	UpdateWithTx(tx *gorm.DB, borrow *models.Borrow) error
	CountByUserID(userID uint) (int64, error)
}

type borrowRepository struct {
	db *gorm.DB
}

func NewBorrowRepository(db *gorm.DB) BorrowRepository {
	return &borrowRepository{db:db}
}

func (r *borrowRepository) Create(borrow *models.Borrow) error {
	return r.db.Create(borrow).Error
}

func (r *borrowRepository) CreateWithTx(tx *gorm.DB,borrow *models.Borrow) error {
	return tx.Create(borrow).Error
}

func (r *borrowRepository) FindByID(id uint) (*models.Borrow, error) {
	var borrow models.Borrow
	err := r.db.Preload("User").Preload("Book").First(&borrow, id).Error
	if err != nil {
		return nil, err
	}
	return &borrow, nil
}

func (r *borrowRepository) FindByIDWithLock(tx *gorm.DB, id uint) (*models.Borrow, error) {
	var borrow models.Borrow
	err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&borrow, id).Error
	if err != nil {
		return nil, err
	}
	return &borrow, nil
}

func (r *borrowRepository) FindByUserID(userID uint, limit, offset int) ([]models.Borrow, error) {
	var borrows []models.Borrow
	err := r.db.Where("user_id = ?", userID).
		Preload("Book").
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&borrows).Error
	return borrows, err
}

func (r *borrowRepository) Update(borrow *models.Borrow) error {
	return r.db.Save(borrow).Error
}

func (r *borrowRepository) UpdateWithTx(tx *gorm.DB, borrow *models.Borrow) error {
	return tx.Save(borrow).Error
}

func (r *borrowRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Borrow{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}