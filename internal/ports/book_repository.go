package ports

import (
	"context"

	"example.com/practice/fiber/internal/domain"
)

type BookRepository interface {
	CreateBook(ctx context.Context, book *domain.Book) error
	GetBookByID(ctx context.Context, id int) (*domain.Book, error)
	GetAllBooks(ctx context.Context) ([]*domain.Book, error)
	UpdateBook(ctx context.Context, book *domain.Book) error
	DeleteBook(ctx context.Context, id int) error
}