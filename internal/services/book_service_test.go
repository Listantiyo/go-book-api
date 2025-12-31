package services

import (
	"book-api/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock BookRepository
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) FindAll(limit, offset int) ([]models.Book, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookRepository) FindByID(id uint) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), nil
}

func (m *MockBookRepository) FindByIDWithLock(tx *gorm.DB, id uint) (*models.Book, error) {
	args := m.Called(tx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), nil
}

func (m *MockBookRepository) FindByISBN(isbn string) (*models.Book, error) {
	args := m.Called(isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) Update(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) UpdateWithTx(tx *gorm.DB, book *models.Book) error {
	args := m.Called(tx, book)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// Test CreateBook - Success
func TestCreateBook_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	// Setup mock
	mockRepo.On("FindByISBN", "123456").Return(nil, errors.New("Not Found"))
	mockRepo.On("Create", mock.AnythingOfType("*models.Book")).Return(nil)

	// Execute
	book, err := service.CreateBook("Test Book", "Test Author", "123456", "Description", 10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "Test Book", book.Title)
	assert.Equal(t, 10, book.Stock)
	mockRepo.AssertExpectations(t)
}

// Test CreateBook - ISBN Already Exists
func TestCreateBook_ISBNAlreadyExists(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	existingBook := &models.Book{
		ID: 1,
		ISBN: "123456",
	}

	// Setup mock
	mockRepo.On("FindByISBN", "123456").Return(existingBook, nil)

	// Execute
	book, err := service.CreateBook("Test Book", "Test Author", "123456", "Description", 10)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, "book with this ISBN already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test CreateBook - Negative Stock
func TestCreateBook_NegativeStock(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	// Execute dengan stock negatif
	book, err := service.CreateBook("Test Book", "Test Author", "123456", "Description", -5)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, "stock cannot be nagtive", err.Error())
	// Tidak perlu cek mock karena validation gagal sebelum hit repository
}

// Test GetAllBooks - Success
func TestGetAllBooks_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockBooks := []models.Book{
		{ID: 1, Title: "Book 1"},
		{ID: 2, Title: "Book 2"},
	}

	// Setup mock
	mockRepo.On("FindAll", 10, 0).Return(mockBooks, nil)
	mockRepo.On("Count").Return(int64(2), nil)

	// Execute
	books, total, err := service.GetAllBooks(0, 10)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(books))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// Test GetBookByID - Success
func TestGetBookByID_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockBook := &models.Book{
		ID: 1,
		Title: "Test Book",
	}

	// Setup mock
	mockRepo.On("FindByID", uint(1)).Return(mockBook, nil)

	// Execute
	book, err := service.GetBookByID(uint(1))

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "Test Book", book.Title)
	mockRepo.AssertExpectations(t)
}

// Test GetBookByID - Not Found
func TestGetBookByID_NotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	// Setup mock
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("book not found"))

	// Execute
	book, err := service.GetBookByID(uint(999))

	// Assert
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, "book not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test DeleteBook - Success
func TestDeleteBook_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockBook := &models.Book{
		ID: 1,
		Title: "Test Book",
	}

	// Setup mock
	mockRepo.On("FindByID", uint(1)).Return(mockBook, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	// Execute
	err := service.DeleteBook(uint(1))

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}