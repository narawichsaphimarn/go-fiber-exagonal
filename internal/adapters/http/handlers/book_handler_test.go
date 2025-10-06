package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"example.com/practice/fiber/internal/adapters/http/middleware"
	domain "example.com/practice/fiber/internal/domain"
	"example.com/practice/fiber/internal/ports"
	usecase "example.com/practice/fiber/internal/usecase/book"
	"github.com/gofiber/fiber/v2"
)

// Mock BookRepository implementing ports.BookRepository
type mockBookRepo struct {
	createErr     error
	getByIDResult *domain.Book
	getByIDErr    error
	getAllResult  []*domain.Book
	getAllErr     error
	updateErr     error
	deleteErr     error
	mu            sync.Mutex
}

func (m *mockBookRepo) CreateBook(ctx context.Context, book *domain.Book) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// simulate ID assignment
	if book.ID == 0 {
		book.ID = 1
	}
	return nil
}
func (m *mockBookRepo) GetBookByID(ctx context.Context, id int) (*domain.Book, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.getByIDResult, nil
}
func (m *mockBookRepo) GetAllBooks(ctx context.Context) ([]*domain.Book, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.getAllResult, nil
}
func (m *mockBookRepo) UpdateBook(ctx context.Context, book *domain.Book) error {
	return m.updateErr
}
func (m *mockBookRepo) DeleteBook(ctx context.Context, id int) error {
	return m.deleteErr
}

// Helper to build protected app for books
func buildBookApp(repo ports.BookRepository, tp ports.TokenProvider) *fiber.App {
	app := fiber.New()
	svc := usecase.NewBookService(repo)
	h := NewBookHandler(svc)
	v1Protected := app.Group("/v1/auth", middleware.NewAuthMiddleware(tp).Protect)
	h.RegisterRoutes(v1Protected)
	return app
}

// Mock TokenProvider (renamed to avoid package-level name collision)
type mockBookTP struct {
	validateUserID string
	validateErr    error
}

func (m *mockBookTP) GenerateToken(userId string) (string, error) { return "", nil }
func (m *mockBookTP) ValidateToken(token string) (string, error) {
	if m.validateErr != nil {
		return "", m.validateErr
	}
	return m.validateUserID, nil
}

// --- Auth middleware cases ---
func TestBooks_Auth_MissingToken(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "1"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/books", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestBooks_Auth_InvalidToken(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateErr: fiber.ErrUnauthorized}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/books", nil)
	req.Header.Set("Authorization", "Bearer bad")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

// --- GET /books ---
func TestGetAllBooks_Success(t *testing.T) {
	repo := &mockBookRepo{getAllResult: []*domain.Book{{ID: 2, Title: "B", Author: "A", Price: 10, Stock: 1}}}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/books", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got []*domain.Book
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if len(got) != 1 || got[0].ID != 2 {
		t.Fatalf("unexpected books: %#v", got)
	}
}

func TestGetAllBooks_RepoError(t *testing.T) {
	repo := &mockBookRepo{getAllErr: fiber.ErrInternalServerError}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/books", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- GET /books/:id ---
func TestGetBookByID_Success(t *testing.T) {
	repo := &mockBookRepo{getByIDResult: &domain.Book{ID: 5, Title: "X", Author: "Y", Price: 1, Stock: 2}}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/books/5", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got domain.Book
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.ID != 5 {
		t.Fatalf("expected id 5, got %d", got.ID)
	}
}

func TestGetBookByID_InvalidID(t *testing.T) {
	repo := &mockBookRepo{getByIDResult: &domain.Book{ID: 5}}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/books/abc", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetBookByID_RepoError(t *testing.T) {
	repo := &mockBookRepo{getByIDErr: fiber.ErrInternalServerError}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/books/5", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- POST /books ---
func TestCreateBook_Success(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	body := `{"title":"T","author":"A","price":10,"stock":1}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/books", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got domain.Book
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.ID == 0 || got.Title != "T" {
		t.Fatalf("unexpected book: %#v", got)
	}
}

func TestCreateBook_InvalidJSON(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/books", bytes.NewBufferString("{"))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateBook_RepoError(t *testing.T) {
	repo := &mockBookRepo{createErr: fiber.ErrInternalServerError}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	body := `{"title":"T","author":"A","price":10,"stock":1}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/books", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- PUT /books/:id ---
func TestUpdateBook_Success(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	body := `{"title":"T2","author":"A","price":11,"stock":2}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/books/7", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got domain.Book
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.ID != 7 || got.Title != "T2" {
		t.Fatalf("unexpected book: %#v", got)
	}
}

func TestUpdateBook_InvalidID(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	body := `{"title":"T2","author":"A","price":11,"stock":2}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/books/xyz", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdateBook_InvalidJSON(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodPut, "/v1/auth/books/7", bytes.NewBufferString("{"))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdateBook_RepoError(t *testing.T) {
	repo := &mockBookRepo{updateErr: fiber.ErrInternalServerError}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	body := `{"title":"T2","author":"A","price":11,"stock":2}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/books/7", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- DELETE /books/:id ---
func TestDeleteBook_Success(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodDelete, "/v1/auth/books/9", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["message"] == "" {
		t.Fatalf("expected delete success message, got %#v", got)
	}
}

func TestDeleteBook_InvalidID(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodDelete, "/v1/auth/books/abc", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestDeleteBook_RepoError(t *testing.T) {
	repo := &mockBookRepo{deleteErr: fiber.ErrInternalServerError}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	req := httptest.NewRequest(http.MethodDelete, "/v1/auth/books/9", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- Concurrency safety ---
func TestCreateBook_Concurrent(t *testing.T) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	body := []byte(`{"title":"T","author":"A","price":10,"stock":1}`)
	const N = 10
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/books", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer good")
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				t.Errorf("concurrent create failed: err=%v status=%d", err, resp.StatusCode)
			}
		}()
	}
	wg.Wait()
}

// --- Benchmarks ---
func BenchmarkGetAllBooks(b *testing.B) {
	repo := &mockBookRepo{getAllResult: []*domain.Book{{ID: 1, Title: "B", Author: "A", Price: 10, Stock: 1}}}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/v1/auth/books", nil)
		req.Header.Set("Authorization", "Bearer good")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			b.Fatalf("unexpected: err=%v status=%d", err, resp.StatusCode)
		}
	}
}

func BenchmarkCreateBook(b *testing.B) {
	repo := &mockBookRepo{}
	tp := &mockBookTP{validateUserID: "3"}
	app := buildBookApp(repo, tp)

	body := []byte(`{"title":"T","author":"A","price":10,"stock":1}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/books", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer good")
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			b.Fatalf("unexpected: err=%v status=%d", err, resp.StatusCode)
		}
	}
}
