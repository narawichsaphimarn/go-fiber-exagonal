package usecase

import (
	"context"

	domains "example.com/practice/fiber/internal/domain"
	"example.com/practice/fiber/internal/ports"
)

type BookService struct {
	repo ports.BookRepository
}

func NewBookService(repo ports.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(ctx context.Context, book *domains.Book) error {
	return s.repo.CreateBook(ctx, book)
}

func (s *BookService) GetBookByID(ctx context.Context, id int) (*domains.Book, error) {
	return s.repo.GetBookByID(ctx, id)
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]*domains.Book, error) {
	return s.repo.GetAllBooks(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, book *domains.Book) error {
	return s.repo.UpdateBook(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.repo.DeleteBook(ctx, id)
}
