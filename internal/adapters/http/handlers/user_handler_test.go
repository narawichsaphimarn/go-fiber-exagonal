package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/practice/fiber/internal/adapters/http/middleware"
	"example.com/practice/fiber/internal/domain"
	"example.com/practice/fiber/internal/ports"
	usercase "example.com/practice/fiber/internal/usecase/user"
	"example.com/practice/fiber/pkg"
	"github.com/gofiber/fiber/v2"
)

// --- Mocks ---
type mockUserRepo struct {
	createUserErr        error
	getUserByIDResult    *domain.User
	getUserByIDErr       error
	getUserByEmailResult *domain.User
	getUserByEmailErr    error
	getAllUsersResult    []*domain.User
	getAllUsersErr       error
	updateUserErr        error
	deleteUserErr        error
	updatePasswordErr    error
}

// Implement ports.UserRepository
func (m *mockUserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	return m.createUserErr
}
func (m *mockUserRepo) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return m.getUserByIDResult, m.getUserByIDErr
}
func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.getUserByEmailResult, m.getUserByEmailErr
}
func (m *mockUserRepo) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return m.getAllUsersResult, m.getAllUsersErr
}
func (m *mockUserRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	return m.updateUserErr
}
func (m *mockUserRepo) DeleteUser(ctx context.Context, id int) error { return m.deleteUserErr }
func (m *mockUserRepo) UpdatePassword(ctx context.Context, id int, password string) error {
	return m.updatePasswordErr
}

type mockTP struct {
	token          string
	validateUserID string
	validateErr    error
}

func (m *mockTP) GenerateToken(userId string) (string, error) { return m.token, nil }
func (m *mockTP) ValidateToken(token string) (string, error) {
	if m.validateErr != nil {
		return "", m.validateErr
	}
	return m.validateUserID, nil
}

// Helper to build app with routes
func buildApp(repo ports.UserRepository, tp ports.TokenProvider) *fiber.App {
	app := fiber.New()
	svc := usercase.NewUserService(repo, tp)
	h := NewUserHandler(svc)

	v1 := app.Group("/v1")
	h.RegisterNotProtected(v1)

	v1Protected := app.Group("/v1/auth", middleware.NewAuthMiddleware(tp).Protect)
	h.RegisterProtected(v1Protected)
	return app
}

// --- RegisterUser ---
func TestRegisterUser_Success(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{}
	app := buildApp(repo, tp)

	body := `{"email":"user@example.com","password":"password123","first_name":"User","last_name":"Example","username":"user123"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var got map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["message"] != "user registered" {
		t.Fatalf("expected message 'user registered', got %v", got)
	}
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{}
	app := buildApp(repo, tp)

	body := `{"email": "bad",` // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestRegisterUser_ValidationError(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{}
	app := buildApp(repo, tp)

	// missing required fields
	body := `{"email":"not-an-email","password":"short","first_name":"","last_name":"","username":"ab"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestRegisterUser_InternalError(t *testing.T) {
	repo := &mockUserRepo{createUserErr: fiber.ErrInternalServerError}
	tp := &mockTP{}
	app := buildApp(repo, tp)

	body := `{"email":"user@example.com","password":"password123","first_name":"User","last_name":"Example","username":"user123"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- LoginUser ---
func TestLoginUser_Success(t *testing.T) {
	hashed, _ := pkg.HashPassword("password123")
	repo := &mockUserRepo{getUserByEmailResult: &domain.User{ID: 3, Email: "user@example.com", Password: hashed, Username: "user123"}}
	tp := &mockTP{token: "token-123"}
	app := buildApp(repo, tp)

	payload := map[string]string{"email": "user@example.com", "password": "password123"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var got map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["token"] != "token-123" {
		t.Fatalf("expected token 'token-123', got %s", got["token"])
	}
}

func TestLoginUser_InvalidJSON(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{}
	app := buildApp(repo, tp)

	body := `{"email":` // malformed
	req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestLoginUser_ValidationError(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{}
	app := buildApp(repo, tp)

	payload := map[string]string{"email": "bad", "password": ""}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	hashed, _ := pkg.HashPassword("password123")
	repo := &mockUserRepo{getUserByEmailResult: &domain.User{ID: 3, Email: "user@example.com", Password: hashed, Username: "user123"}}
	tp := &mockTP{}
	app := buildApp(repo, tp)

	payload := map[string]string{"email": "user@example.com", "password": "wrong"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- Protected: GetUserByID ---
func TestGetUserByID_MissingToken(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/user/3", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestGetUserByID_InvalidToken(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateErr: fiber.ErrUnauthorized}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/user/3", nil)
	req.Header.Set("Authorization", "Bearer bad")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	repo := &mockUserRepo{getUserByIDResult: &domain.User{ID: 3, Email: "user@example.com"}}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/user/3", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var u domain.User
	_ = json.NewDecoder(resp.Body).Decode(&u)
	if u.ID != 3 {
		t.Fatalf("expected user ID 3, got %d", u.ID)
	}
}

func TestGetUserByID_BadParam(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/user/abc", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetUserByID_RepoError(t *testing.T) {
	repo := &mockUserRepo{getUserByIDErr: fiber.ErrInternalServerError}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/user/3", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- Protected: GetUserByEmail ---
func TestGetUserByEmail_Success(t *testing.T) {
	repo := &mockUserRepo{getUserByEmailResult: &domain.User{ID: 3, Email: "user@example.com"}}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/user/email/user@example.com", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGetUserByEmail_RepoError(t *testing.T) {
	repo := &mockUserRepo{getUserByEmailErr: fiber.ErrInternalServerError}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/user/email/user@example.com", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- Protected: GetAllUsers ---
func TestGetAllUsers_Success(t *testing.T) {
	repo := &mockUserRepo{getAllUsersResult: []*domain.User{{ID: 1}, {ID: 2}}}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/users", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var users []domain.User
	_ = json.NewDecoder(resp.Body).Decode(&users)
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestGetAllUsers_RepoError(t *testing.T) {
	repo := &mockUserRepo{getAllUsersErr: fiber.ErrInternalServerError}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/users", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- Protected: UpdateUser ---
func TestUpdateUser_Success(t *testing.T) {
	repo := &mockUserRepo{
		getUserByIDResult: &domain.User{ID: 3, Email: "user@example.com"},
	}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"first_name":"Updated","last_name":"User"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/3", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestUpdateUser_BadParam(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"first_name":"Updated","last_name":"User"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/abc", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"first_name":` // malformed
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/3", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdateUser_ValidationError(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"first_name":"","last_name":""}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/3", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdateUser_RepoError(t *testing.T) {
	repo := &mockUserRepo{
		getUserByIDResult: &domain.User{ID: 3, Email: "user@example.com"},
		updateUserErr:     fiber.ErrInternalServerError,
	}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"first_name":"Updated","last_name":"User"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/3", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- Protected: DeleteUser ---
func TestDeleteUser_Success(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodDelete, "/v1/auth/user/3", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDeleteUser_BadParam(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodDelete, "/v1/auth/user/abc", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestDeleteUser_RepoError(t *testing.T) {
	repo := &mockUserRepo{deleteUserErr: fiber.ErrInternalServerError}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	req := httptest.NewRequest(http.MethodDelete, "/v1/auth/user/3", nil)
	req.Header.Set("Authorization", "Bearer good")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

// --- Protected: UpdatePassword ---
func TestUpdatePassword_Success(t *testing.T) {
	repo := &mockUserRepo{
		getUserByIDResult: &domain.User{ID: 3, Email: "user@example.com"},
	}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"new_password":"newpassword123"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/3/password", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestUpdatePassword_BadParam(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"new_password":"newpassword123"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/abc/password", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdatePassword_InvalidJSON(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"new_password":` // malformed
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/3/password", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdatePassword_ValidationError(t *testing.T) {
	repo := &mockUserRepo{}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"new_password":"short"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/3/password", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdatePassword_RepoError(t *testing.T) {
	repo := &mockUserRepo{updatePasswordErr: fiber.ErrInternalServerError, getUserByIDResult: &domain.User{ID: 3, Email: "user@example.com"}}
	tp := &mockTP{validateUserID: "3"}
	app := buildApp(repo, tp)

	body := `{"new_password":"newpassword123"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/auth/user/3/password", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer good")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}
