package services

import (
	"book-api/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}
// Create
func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
// FindByEmail
func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
// FindByID
func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(0)
	}
	return args.Get(0).(*models.User), nil
}

// Test Register - Success
func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Setup mock expectation
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Execute
	user, err := service.Register("Test User", "test@example.com", "password123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, user.Password)
	mockRepo.AssertExpectations(t)
}

// Test Register - Email Already Exist
func TestRegister_EmailAlreadyExist(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &models.User{
		ID: 1,
		Email: "test@example.com",
	}

	// Setup mock -  email sudah ada
	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, errors.New("invalid email or password"))

	// Execute
	user, err := service.Register("Test User", "test@example.com", "password123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already registered", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test Login - Success
func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Buat user dengan password yang sudah di-hash
	// Password asli: "password123"
	hashedPassword := "$2a$12$Vobb3BoaYxoJKIwDkGX7kuqNSs/Jr61HdBR7GEr5yD.OhMJBDzAmS"
	existingUser := &models.User{
		ID: 1,
		Email: "test@example.com",
		Password: hashedPassword,
	}

	// Setup mock
	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	// Execute
	token, err := service.Login("test@example.com", "password123", "secret-key")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

// Test Login - Invalid Password
func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	hashedPassword := "$2a$12$Vobb3BoaYxoJKIwDkGX7kuqNSs/Jr61HdBR7GEr5yD.OhMJBDzAmS"
	existingUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	// Setup mock
	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	// Execute dengan password salah
	token, err := service.Login("test@example.com", "wrongpassword", "secret-key")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test Login - User Not Found
func TestLogin_UserNotFoud(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Setup mock - user tidak ditemukan
	mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("not found"))

	// Execute
	token, err := service.Login("notfound@example.com", "password123", "secret-key")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}