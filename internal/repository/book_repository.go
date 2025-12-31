package repository

import (
	"book-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookRepository interface {
	Create(book *models.Book) error
	FindAll(limit, offset int) ([]models.Book, error)
	FindByID(id uint) (*models.Book, error)
	FindByIDWithLock(tx *gorm.DB, id uint) (*models.Book, error)
	FindByISBN(isbn string) (*models.Book, error)
	Update(book *models.Book) error
	UpdateWithTx(tx *gorm.DB, book *models.Book) error
	Delete(id uint) error
	Count() (int64, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) FindAll(limit, offset int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Limit(limit).Offset(offset).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) FindByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.db.Where("id = ?", id).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}
func (r *bookRepository) FindByIDWithLock(tx *gorm.DB, id uint) (*models.Book, error) {
	var book models.Book
	err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) FindByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	err := r.db.Where("isbn = ?", isbn).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) UpdateWithTx(tx *gorm.DB, book *models.Book) error {
	return tx.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}

func (r *bookRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Book{}).Count(&count).Error
	return count, err
}