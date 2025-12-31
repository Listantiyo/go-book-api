package services

import (
	"book-api/internal/models"
	"book-api/internal/repository"
	"errors"
)

type BookService interface {
	CreateBook(title, author, isbn, description string, stock int) (*models.Book, error)
	GetAllBooks(page, pageSize int) ([]models.Book, int64, error)
	GetBookByID(id uint) (*models.Book, error)
	UpdateBook(id uint, title, author, isbn, description string, stock int) (*models.Book, error)
	DeleteBook(id uint) error
}

type bookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{bookRepo: bookRepo}
}

func (s *bookService) CreateBook(title, author, isbn, description string, stock int) (*models.Book, error) {
	// Validasi stock tidak boleh negatif
	if stock < 0 {
		return nil, errors.New("stock cannot be nagtive")
	}

	// Cek apakah ISBN sudah ada
	existingBook, _ := s.bookRepo.FindByISBN(isbn)
	if existingBook != nil {
		return nil, errors.New("book with this ISBN already exists")
	}

	newBook := models.Book{
		Title: title,
		Author: author,
		ISBN: isbn,
		Description: description,
		Stock: stock,
	}

	if err := s.bookRepo.Create(&newBook); err != nil {
		return nil, err
	}

	return &newBook, nil
}

func (s *bookService) GetAllBooks(page, pageSize int) ([]models.Book, int64, error) {
	// Default pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	books, err := s.bookRepo.FindAll(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.bookRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (s *bookService) GetBookByID(id uint) (*models.Book, error) {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *bookService) UpdateBook(id uint, title, author, isbn, description string, stock int) (*models.Book, error) {
	// Cek apakah buku ada
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Validasi stock tidak boleh negatif
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	// Update fileds
	book.Title = title
	book.Author = author
	book.ISBN = isbn
	book.Description = description
	book.Stock = stock

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) DeleteBook(id uint) error {
	// Cek apakah buku ada
	if _, err := s.bookRepo.FindByID(id); err != nil {
		return err
	}
	
	return s.bookRepo.Delete(id)
}