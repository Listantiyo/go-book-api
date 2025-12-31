package services

import (
	"book-api/internal/database"
	"book-api/internal/models"
	"book-api/internal/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

type BorrowService interface {
	BorrowBook(userID, bookID uint) (*models.Borrow, error)
	ReturnBook(borrowID uint) (*models.Borrow, error)
	GetUserBorrows(userID uint, page, pageSize int) ([]models.Borrow, int64, error)
	GetBorrowByID(borrowID uint) (*models.Borrow, error)
}

type borrowService struct {
	borrowRepo 	repository.BorrowRepository
	bookRepo 	repository.BookRepository
	txManager 	database.TransactionManager
}

func NewBorrowService(
	borrowRepo repository.BorrowRepository,
	bookRepo repository.BookRepository,
	txManager database.TransactionManager,
) BorrowService {
	return &borrowService{
		borrowRepo: borrowRepo,
		bookRepo: 	bookRepo,
		txManager:	txManager,
	}
}

func (s *borrowService) BorrowBook(userID, bookID uint) (*models.Borrow, error) {
	var result *models.Borrow

	// Semua operasi dalam transaction
	err := s.txManager.WithTransaction(func(tx *gorm.DB) error {
		// 1. Cek dan LOCK buku
		book, err := s.bookRepo.FindByIDWithLock(tx, bookID)
		if err != nil {
			return errors.New("book not found")
		}

		if book.Stock <= 0 {
			return errors.New("book out of stock")
		}

		// 2. Kurangi stock buku
		book.Stock--
		if err := s.bookRepo.UpdateWithTx(tx, book); err != nil {
			return err
		}

		// 3. Buat record borrow
		borrow := &models.Borrow{
			UserID: userID,
			BookID: bookID,
			BorrowDate: time.Now(),
			DueDate: time.Now().Add(14 * 24 * time.Hour), // dua minggu
			Status: models.BorrowStatusBorrowed,
		}

		if err := s.borrowRepo.CreateWithTx(tx, borrow); err != nil {
			return err
		}

		result = borrow
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *borrowService) ReturnBook(borrowID uint) (*models.Borrow, error) {
	var result *models.Borrow

	err := s.txManager.WithTransaction(func(tx *gorm.DB) error {
		// 1. Cari dan LOCK borrow record
		borrow, err := s.borrowRepo.FindByIDWithLock(tx, borrowID)
		if err != nil {
			return errors.New("borrow record not found")
		}
		// 2. Cek apakah sudah dikembalikan
		if borrow.Status == models.BorrowStatusReturned {
			return errors.New("book already returned")
		}
		// 3. Update status dan return date
		now := time.Now()
		borrow.ReturnDate = &now
		borrow.Status = models.BorrowStatusReturned
		if err := s.borrowRepo.UpdateWithTx(tx, borrow); err != nil {
			return err
		}
		// 4. Lock dan Update stock buku
		book, err := s.bookRepo.FindByIDWithLock(tx, borrow.BookID); 
		if err != nil {
			return err
		}

		book.Stock++
		if err := s.bookRepo.UpdateWithTx(tx, book); err != nil {
			return err
		}

		result = borrow
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *borrowService) GetUserBorrows(userID uint, page, pageSize int) ([]models.Borrow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	borrows, err := s.borrowRepo.FindByUserID(userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.borrowRepo.CountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	return borrows, total, nil
}

func (s *borrowService) GetBorrowByID(borrowID uint) (*models.Borrow, error) {
	borrow, err := s.borrowRepo.FindByID(borrowID)
	if err != nil {
		return nil, errors.New("borrow record not found")
	}
	return borrow, nil
}