package services

import (
	"book-api/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockBorrowRepository
type MockBorrowRepository struct {
	mock.Mock
}

func (m *MockBorrowRepository) Create(borrow *models.Borrow) error {
	args := m.Called(borrow)
	return args.Error(0)
}
func (m *MockBorrowRepository) CreateWithTx(tx *gorm.DB, borrow *models.Borrow) error {
	args := m.Called(tx, borrow)
	return args.Error(0)
}
func (m *MockBorrowRepository) FindByID(id uint) (*models.Borrow, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Borrow), nil
}
func (m *MockBorrowRepository) FindByIDWithLock(tx *gorm.DB, id uint) (*models.Borrow, error) {
	args := m.Called(tx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Borrow), nil
}
func (m *MockBorrowRepository) FindByUserID(userID uint, limit, offset int) ([]models.Borrow, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]models.Borrow), nil
}
func (m *MockBorrowRepository) Update(borrow *models.Borrow) error {
	args := m.Called(borrow)
	return args.Error(0)
}
func (m *MockBorrowRepository) UpdateWithTx(tx *gorm.DB, borrow *models.Borrow) error {
	args := m.Called(tx, borrow)
	return args.Error(0)
}
func (m *MockBorrowRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)	
}

// MockTransactionManager
type MockTransactionManager struct {
	mock.Mock
}
func (m *MockTransactionManager) WithTransaction(fn func(*gorm.DB) error) error {
	return fn(nil)
}

// TestBorrowBook - Success
func TestBorrowBook_Success(t *testing.T) {
	mockBorrowRepo 	:= new(MockBorrowRepository)
	mockBookRepo 	:= new(MockBookRepository)
	mockTxManager	:= new(MockTransactionManager)
	service := NewBorrowService(mockBorrowRepo, mockBookRepo, mockTxManager)

	book := &models.Book{
		ID: 2,
		Title: "Test Book",
		Stock: 2,
	}
	// Buat ekspektasi
	mockBookRepo.On("FindByIDWithLock", mock.Anything, uint(2)).Return(book, nil)
	mockBookRepo.On("UpdateWithTx", mock.Anything, mock.AnythingOfType("*models.Book")).Return(nil)
	mockBorrowRepo.On("CreateWithTx", mock.Anything, mock.AnythingOfType("*models.Borrow")).Return(nil)

	// Execute
	borrow, err := service.BorrowBook(uint(2), uint(2))

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, borrow)
	assert.Equal(t, uint(2), borrow.UserID)
	assert.Equal(t, uint(2), borrow.BookID)
	assert.Equal(t, borrow.Status, models.BorrowStatusBorrowed)
	mockBorrowRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

// TestBorrowBook - Out of Stock
func TestBorrowBook_OutOfStock(t *testing.T) {
	mockBorrowRepo := new(MockBorrowRepository)
	mockBookRepo := new(MockBookRepository)
	mockTxManager := new(MockTransactionManager)
	service := NewBorrowService(mockBorrowRepo, mockBookRepo, mockTxManager)

	book := &models.Book{
		ID: 2,
		Title: "Test Book",
		Stock: 0,
	}

	// Expectations
	mockBookRepo.On("FindByIDWithLock", mock.Anything, uint(2)).Return(book, nil)

	// Execute
	borrow, err := service.BorrowBook(uint(2), uint(2))

	assert.Error(t, err)
	assert.Nil(t, borrow)
	assert.Equal(t, err.Error(), "book out of stock")
	mockBookRepo.AssertExpectations(t)
}

// TestBorrowBook - Book Not Found
func TestBorrowBook_BookNotFound(t *testing.T) {
	mockBorrowRepo := new(MockBorrowRepository)
	mockBookRepo := new(MockBookRepository)
	mockTxManager := new(MockTransactionManager)
	service := NewBorrowService(mockBorrowRepo, mockBookRepo, mockTxManager)

	// Expectations
	mockBookRepo.On("FindByIDWithLock", mock.Anything, uint(999)).Return(nil, errors.New("not found"))

	// Execute
	borrow, err := service.BorrowBook(1, uint(999))

	assert.Error(t, err)
	assert.Nil(t, borrow)
	assert.Equal(t, err.Error(), "book not found")
	mockBookRepo.AssertExpectations(t)
}

// TestReturnBook - Success
func TestReturnBook_Success(t *testing.T) {
	mockBorrowRepo := new(MockBorrowRepository)
	mockBookRepo := new(MockBookRepository)
	mockTxManager := new(MockTransactionManager)
	service := NewBorrowService(mockBorrowRepo, mockBookRepo, mockTxManager)

	borrow := &models.Borrow{
		ID: 1,
		UserID: 1,
		BookID: 1,
		Status: models.BorrowStatusBorrowed,
	}
	book := &models.Book{
		ID: 1,
		Title: "Test Book",
		Stock: 2,
	}
	// Expectations
	mockBorrowRepo.On("FindByIDWithLock", mock.Anything, uint(1)).Return(borrow, nil)
	mockBorrowRepo.On("UpdateWithTx", mock.Anything, mock.AnythingOfType("*models.Borrow")).Return(nil)
	mockBookRepo.On("FindByIDWithLock", mock.Anything, uint(1)).Return(book, nil)
	mockBookRepo.On("UpdateWithTx", mock.Anything, mock.AnythingOfType("*models.Book")).Return(nil)

	// Execute
	borrow, err := service.ReturnBook(uint(1))

	// Asserts
	assert.NoError(t, err)
	assert.NotNil(t, borrow)
	assert.Equal(t, borrow.Status, models.BorrowStatusReturned)
	assert.NotNil(t, borrow.ReturnDate)
	mockBorrowRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}