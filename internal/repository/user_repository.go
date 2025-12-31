package repository

import (
	"book-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
// Implement method Create
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
// Implement method FindByEmail
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return  nil, err
	}

	return &user, nil
}
// Implement method FindByID
func (r *userRepository) FindByID(id uint) (*models.User, error){
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}