package services

import (
	"book-api/internal/models"
	"book-api/internal/repository"
	"book-api/internal/utils"
	"errors"
)

type AuthService interface {
	Register(name, email, password string) (*models.User, error)
	Login(email, password, jwtSecret string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(name, email, password string) (*models.User, error) {
	// Validasi cek email sudah terdaftar
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Buat user baru
	newUser := models.User{
		Name: name,
		Email: email,
		Password: hashedPassword,
	}

	//Simpan user ke repository
	if err := s.userRepo.Create(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *authService) Login(email, password, jwtSecret string) (string, error) {
	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Cek password
	if !utils.CheckHashPassword(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	// Generate token JWT
	token, err := utils.GenerateToken(user.ID, user.Email, jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}