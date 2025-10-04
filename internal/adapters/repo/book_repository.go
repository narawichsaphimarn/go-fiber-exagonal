package repo

import (
	"context"
	"time"

	domain "example.com/practice/fiber/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepo struct {
	db *pgxpool.Pool
}


func NewBookRepo(db *pgxpool.Pool) *BookRepo {
	return &BookRepo{db: db}
}

func (r *BookRepo) CreateBook(ctx context.Context, book *domain.Book) error {
	_, err := r.db.Exec(ctx, "INSERT INTO books (title, author, price, stock) VALUES ($1, $2, $3, $4)", book.Title, book.Author, book.Price, book.Stock)
	return err
}

func (r *BookRepo) GetBookByID(ctx context.Context, id int) (*domain.Book, error) {
	var book domain.Book
	err := r.db.QueryRow(ctx, "SELECT id, title, author, price, stock, created_at, updated_at FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Stock, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepo) GetAllBooks(ctx context.Context) ([]*domain.Book, error) {
	rows, err := r.db.Query(ctx, "SELECT id, title, author, price, stock, created_at, updated_at FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book		
	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Stock, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookRepo) UpdateBook(ctx context.Context, book *domain.Book) error {
	_, err := r.db.Exec(ctx, "UPDATE books SET title = $1, author = $2, price = $3, stock = $4, updated_at = $5 WHERE id = $6", book.Title, book.Author, book.Price, book.Stock, time.Now(), book.ID)
	return err
}

func (r *BookRepo) DeleteBook(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM books WHERE id = $1", id)
	return err
}