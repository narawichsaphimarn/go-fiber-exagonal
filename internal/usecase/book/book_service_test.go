package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"example.com/practice/fiber/internal/domain"
)

type mockBookRepository struct {
	books domain.Book
	err   error
}

func (m *mockBookRepository) CreateBook(ctx context.Context, book *domain.Book) error {
	return m.err
}

func (m *mockBookRepository) GetBookByID(ctx context.Context, id int) (*domain.Book, error) {
	return &m.books, m.err
}

func (m *mockBookRepository) GetAllBooks(ctx context.Context) ([]*domain.Book, error) {
	if m.err != nil {
		return []*domain.Book{}, m.err
	}
	return []*domain.Book{&m.books}, nil
}

func (m *mockBookRepository) UpdateBook(ctx context.Context, book *domain.Book) error {
	return m.err
}

func (m *mockBookRepository) DeleteBook(ctx context.Context, id int) error {
	return m.err
}

func TestCreateBook(t *testing.T) {
	repo := &mockBookRepository{}
	service := NewBookService(repo)

	book := &domain.Book{
		Title:     "Test Book",
		Author:    "Test Author",
		Price:     100,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := service.CreateBook(context.Background(), book)
	if err != nil {
		t.Errorf("CreateBook() error = %v, wantErr %v", err, repo.err)
	}
}

func TestCreateBook_Error(t *testing.T) {
	repo := &mockBookRepository{err: errors.New("test error")}
	service := NewBookService(repo)

	book := &domain.Book{
		Title:     "Test Book",
		Author:    "Test Author",
		Price:     100,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := service.CreateBook(context.Background(), book)
	if err == nil {
		t.Errorf("CreateBook() error = %v, wantErr %v", err, repo.err)
	}
}

func TestGetBookByID(t *testing.T) {
	repo := &mockBookRepository{books: domain.Book{
		ID:        1,
		Title:     "Test Book",
		Author:    "Test Author",
		Price:     100,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}}
	service := NewBookService(repo)

	book, err := service.GetBookByID(context.Background(), 1)
	if err != nil {
		t.Errorf("GetBookByID() error = %v, wantErr %v", err, repo.err)
	}
	if book.ID != repo.books.ID {
		t.Errorf("GetBookByID() book.ID = %v, want %v", book.ID, repo.books.ID)
	}
}

func TestGetBookByID_Error(t *testing.T) {
	repo := &mockBookRepository{books: domain.Book{}, err: errors.New("test error")}
	service := NewBookService(repo)

	book, err := service.GetBookByID(context.Background(), 1)
	fmt.Println(book.ID)
	if err == nil {
		t.Errorf("GetBookByID() error = %v, wantErr %v", err, repo.err)
	}
	if book.ID != 0 {
		t.Errorf("GetBookByID() book = %v, want nil", book)
	}
}

func TestGetAllBooks(t *testing.T) {
	repo := &mockBookRepository{books: domain.Book{
		ID:        1,
		Title:     "Test Book",
		Author:    "Test Author",
		Price:     100,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}}
	service := NewBookService(repo)

	books, err := service.GetAllBooks(context.Background())
	if err != nil {
		t.Errorf("GetAllBooks() error = %v, wantErr %v", err, repo.err)
	}
	if len(books) != 1 {
		t.Errorf("GetAllBooks() len(books) = %v, want %v", len(books), 1)
	}
	if books[0].ID != repo.books.ID {
		t.Errorf("GetAllBooks() books[0].ID = %v, want %v", books[0].ID, repo.books.ID)
	}
}

func TestGetAllBooks_Error(t *testing.T) {
	repo := &mockBookRepository{books: domain.Book{}, err: errors.New("test error")}
	service := NewBookService(repo)

	books, err := service.GetAllBooks(context.Background())
	if err == nil {
		t.Errorf("GetAllBooks() error = %v, wantErr %v", err, repo.err)
	}
	if len(books) != 0 {
		t.Errorf("GetAllBooks() len(books) = %v, want %v", len(books), 0)
	}
}

func TestUpdateBook(t *testing.T) {
	repo := &mockBookRepository{}
	service := NewBookService(repo)

	book := &domain.Book{
		ID:        1,
		Title:     "Test Book",
		Author:    "Test Author",
		Price:     100,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := service.UpdateBook(context.Background(), book)
	if err != nil {
		t.Errorf("UpdateBook() error = %v, wantErr %v", err, repo.err)
	}
}

func TestUpdateBook_Error(t *testing.T) {
	repo := &mockBookRepository{err: errors.New("test error")}
	service := NewBookService(repo)

	book := &domain.Book{
		ID:        1,
		Title:     "Test Book",
		Author:    "Test Author",
		Price:     100,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := service.UpdateBook(context.Background(), book)
	if err == nil {
		t.Errorf("UpdateBook() error = %v, wantErr %v", err, repo.err)
	}
}

func TestDeleteBook(t *testing.T) {
	repo := &mockBookRepository{}
	service := NewBookService(repo)

	err := service.DeleteBook(context.Background(), 1)
	if err != nil {
		t.Errorf("DeleteBook() error = %v, wantErr %v", err, repo.err)
	}
}

func TestDeleteBook_Error(t *testing.T) {
	repo := &mockBookRepository{err: errors.New("test error")}
	service := NewBookService(repo)

	err := service.DeleteBook(context.Background(), 1)
	if err == nil {
		t.Errorf("DeleteBook() error = %v, wantErr %v", err, repo.err)
	}
}
