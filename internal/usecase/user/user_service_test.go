package usercase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"example.com/practice/fiber/internal/domain"
	"example.com/practice/fiber/pkg"
)

type mockUserRepository struct {
	users []domain.User
	err   error
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return m.err
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if len(m.users) == 0 && m.err == nil {
		return nil, nil
	}
	if m.err != nil {
		return nil, m.err
	}
	return &m.users[0], m.err
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if len(m.users) == 0 && m.err == nil {
		return nil, nil
	}
	if m.err != nil {
		return nil, m.err
	}
	return &m.users[0], nil
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	if len(m.users) == 0 && m.err == nil {
		return []*domain.User{}, nil
	}
	if m.err != nil {
		return nil, m.err
	}
	return []*domain.User{&m.users[0]}, m.err
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return m.err
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, id int) error {
	return m.err
}

func (m *mockUserRepository) UpdatePassword(ctx context.Context, id int, password string) error {
	return m.err
}

type mockTokenProvider struct {
	token string
	sub   string
	err   error
}

func (m *mockTokenProvider) GenerateToken(userId string) (string, error) {
	return m.token, m.err
}

func (m *mockTokenProvider) ValidateToken(tokenString string) (string, error) {
	return m.sub, m.err
}

func TestUserService_Login(t *testing.T) {
	// Test Login with valid credentials
	passwordHash, _ := pkg.HashPassword("password123")
	repo := &mockUserRepository{
		users: []domain.User{
			{
				ID:       1,
				Email:    "test@example.com",
				Password: passwordHash,
			},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login
	token, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Errorf("Login failed: %v", err)
	}
	if token != tokenProvider.token {
		t.Errorf("Login token mismatch: expected %s, got %s", tokenProvider.token, token)
	}
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	// Test Login with invalid credentials
	passwordHash, _ := pkg.HashPassword("password123")
	repo := &mockUserRepository{
		users: []domain.User{
			{
				ID:       1,
				Email:    "test@example.com",
				Password: passwordHash,
			},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login with invalid credentials
	_, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Errorf("Login should have failed with invalid credentials")
	}
}

func TestUserService_Login_InvalidEmail(t *testing.T) {
	// Test Login with invalid email
	repo := &mockUserRepository{
		users: []domain.User{},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login with invalid email
	_, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err == nil || err != pkg.ErrUserNotFound {
		t.Errorf("Login should have failed with invalid email: %v", err)
	}
}

func TestUserService_Login_EmptyEmail(t *testing.T) {
	// Test Login with empty email
	repo := &mockUserRepository{
		users: []domain.User{},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login with empty email
	_, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "",
		Password: "password123",
	})
	if err == nil {
		t.Errorf("Login should have failed with empty email: %v", err)
	}
}

func TestUserService_Login_EmptyPassword(t *testing.T) {
	// Test Login with empty password
	repo := &mockUserRepository{
		users: []domain.User{},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login with empty password
	_, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "test@example.com",
		Password: "",
	})
	if err == nil {
		t.Errorf("Login should have failed with empty password: %v", err)
	}
}

func TestUserService_Login_InvalidEmailFormat(t *testing.T) {
	// Test Login with invalid email format
	repo := &mockUserRepository{
		users: []domain.User{},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login with invalid email format
	_, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "testexample.com",
		Password: "password123",
	})
	if err == nil {
		t.Errorf("Login should have failed with invalid email format: %v", err)
	}
}

func TestUserService_Login_EmptyEmailAndPassword(t *testing.T) {
	// Test Login with empty email and password
	repo := &mockUserRepository{
		users: []domain.User{},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login with empty email and password
	_, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "",
		Password: "",
	})
	if err == nil {
		t.Errorf("Login should have failed with empty email and password: %v", err)
	}
}

func TestUserService_Login_GenerateToken(t *testing.T) {
	// Test Login with valid credentials
	passwordHash, _ := pkg.HashPassword("password123")
	repo := &mockUserRepository{
		users: []domain.User{
			{
				ID:       1,
				Email:    "test@example.com",
				Password: passwordHash,
			},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login
	token, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Errorf("Login failed: %v", err)
	}
	if token != tokenProvider.token {
		t.Errorf("Login token mismatch: expected %s, got %s", tokenProvider.token, token)
	}
}

func TestUserService_Login_GenerateToken_Error(t *testing.T) {
	// Test Login with invalid email
	repo := &mockUserRepository{
		users: []domain.User{},
	}
	tokenProvider := &mockTokenProvider{
		err: errors.New("token generation error"),
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Login with invalid email
	_, err := userService.Login(context.Background(), &domain.AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Errorf("Login should have failed with invalid email: %v", err)
	}
}

func TestUserService_Register(t *testing.T) {
	// Test Register with valid user
	repo := &mockUserRepository{}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Register with valid user
	err := userService.Register(context.Background(), &domain.User{
		Email:     "test@example.com",
		Password:  "password123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	})
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}
}

func TestUserService_Register_EmptyEmail(t *testing.T) {
	// Test Register with empty email
	repo := &mockUserRepository{}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Register with empty email
	err := userService.Register(context.Background(), &domain.User{
		Email:     "",
		Password:  "password123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	})
	if err == nil {
		t.Errorf("Register should have failed with empty email: %v", err)
	}
}

func TestUserService_Register_EmptyPassword(t *testing.T) {
	// Test Register with empty password
	repo := &mockUserRepository{}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Register with empty password
	err := userService.Register(context.Background(), &domain.User{
		Email:     "test@example.com",
		Password:  "",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	})
	if err == nil {
		t.Errorf("Register should have failed with empty password: %v", err)
	}
}

func TestUserService_Register_EmptyUsername(t *testing.T) {
	// Test Register with empty username
	repo := &mockUserRepository{}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Register with empty username
	err := userService.Register(context.Background(), &domain.User{
		Email:     "test@example.com",
		Password:  "password123",
		Username:  "",
		FirstName: "Test",
		LastName:  "User",
	})
	if err == nil {
		t.Errorf("Register should have failed with empty username: %v", err)
	}
}

func TestUserService_Register_EmptyFirstName(t *testing.T) {
	// Test Register with empty first name
	repo := &mockUserRepository{}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Register with empty first name
	err := userService.Register(context.Background(), &domain.User{
		Email:     "test@example.com",
		Password:  "password123",
		Username:  "testuser",
		FirstName: "",
		LastName:  "User",
	})
	if err == nil {
		t.Errorf("Register should have failed with empty first name: %v", err)
	}
}

func TestUserService_Register_EmptyLastName(t *testing.T) {
	// Test Register with empty last name
	repo := &mockUserRepository{}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Register with empty last name
	err := userService.Register(context.Background(), &domain.User{
		Email:     "test@example.com",
		Password:  "password123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "",
	})
	if err == nil {
		t.Errorf("Register should have failed with empty last name: %v", err)
	}
}

func TestUserService_Register_UserAlreadyExists(t *testing.T) {
	// Test Register with user already exists
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test Register with user already exists
	err := userService.Register(context.Background(), &domain.User{
		Email:     "test@example.com",
		Password:  "password123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	})
	if err != pkg.ErrUserAlreadyExists {
		t.Errorf("Register should have failed with user already exists: %v", err)
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	// Test UpdatePassword with valid user
	password, err := pkg.HashPassword("newpassword123")
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com", Password: password},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdatePassword with valid user
	err = userService.UpdatePassword(context.Background(), 1, domain.UpdatePasswordRequest{
		NewPassword: "newpassword123",
	})
	if err != nil {
		t.Errorf("UpdatePassword failed: %v", err)
	}
}

func TestUserService_UpdatePassword_UserNotFound(t *testing.T) {
	// Test UpdatePassword with user not found
	repo := &mockUserRepository{
		err: pkg.ErrUserNotFound,
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdatePassword with user not found
	err := userService.UpdatePassword(context.Background(), 2, domain.UpdatePasswordRequest{
		NewPassword: "newpassword123",
	})
	if err != pkg.ErrUserNotFound {
		t.Errorf("UpdatePassword should have failed with user not found: %v", err)
	}
}

func TestUserService_UpdatePassword_EmptyPassword(t *testing.T) {
	// Test UpdatePassword with empty password
	repo := &mockUserRepository{}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdatePassword with empty password
	err := userService.UpdatePassword(context.Background(), 1, domain.UpdatePasswordRequest{
		NewPassword: "",
	})
	if err == nil {
		t.Errorf("UpdatePassword should have failed with empty password: %v", err)
	}
}

func TestUserService_UpdatePassword_InvalidPassword(t *testing.T) {
	// Test UpdatePassword with invalid password
	repo := &mockUserRepository{}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdatePassword with invalid password
	err := userService.UpdatePassword(context.Background(), 1, domain.UpdatePasswordRequest{
		NewPassword: "123456",
	})
	if err == nil {
		t.Errorf("UpdatePassword should have failed with invalid password: %v", err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	// Test UpdateUser with valid user
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdateUser with valid user
	err := userService.UpdateUser(context.Background(), 1, &domain.UpdateUserRequest{
		FirstName: "New",
		LastName:  "User",
	})
	if err != nil {
		t.Errorf("UpdateUser failed: %v", err)
	}
}

func TestUserService_UpdateUser_UserNotFound(t *testing.T) {
	// Test UpdateUser with user not found
	repo := &mockUserRepository{
		err: pkg.ErrUserNotFound,
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdateUser with user not found
	err := userService.UpdateUser(context.Background(), 2, &domain.UpdateUserRequest{
		FirstName: "New",
		LastName:  "User",
	})
	if err != pkg.ErrUserNotFound {
		t.Errorf("UpdateUser should have failed with user not found: %v", err)
	}
}

func TestUserService_UpdateUser_EmptyFirstName(t *testing.T) {
	// Test UpdateUser with empty first name
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdateUser with empty first name
	err := userService.UpdateUser(context.Background(), 1, &domain.UpdateUserRequest{
		FirstName: "",
		LastName:  "User",
	})
	if err == nil {
		t.Errorf("UpdateUser should have failed with empty first name: %v", err)
	}
}

func TestUserService_UpdateUser_EmptyLastName(t *testing.T) {
	// Test UpdateUser with empty last name
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdateUser with empty last name
	err := userService.UpdateUser(context.Background(), 1, &domain.UpdateUserRequest{
		FirstName: "New",
		LastName:  "",
	})
	if err == nil {
		t.Errorf("UpdateUser should have failed with empty last name: %v", err)
	}
}

func TestUserService_UpdateUser_InvalidUserID(t *testing.T) {
	// Test UpdateUser with invalid user ID
	repo := &mockUserRepository{
		users: []domain.User{},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test UpdateUser with invalid user ID
	err := userService.UpdateUser(context.Background(), 0, &domain.UpdateUserRequest{
		FirstName: "New",
		LastName:  "User",
	})
	if err == nil {
		t.Errorf("UpdateUser should have failed with invalid user ID: %v", err)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	// Test GetUserByID with valid user ID
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test GetUserByID with valid user ID
	user, err := userService.GetUserByID(context.Background(), "1")
	if err != nil {
		t.Errorf("GetUserByID failed: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("GetUserByID returned incorrect user ID: %d", user.ID)
	}
}

func TestUserService_GetUserByID_UserNotFound(t *testing.T) {
	// Test GetUserByID with user not found
	repo := &mockUserRepository{
		err: pkg.ErrUserNotFound,
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test GetUserByID with user not found
	_, err := userService.GetUserByID(context.Background(), "2")
	fmt.Println(err)
	if err != pkg.ErrUserNotFound {
		t.Errorf("GetUserByID should have failed with user not found: %v", err)
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	// Test GetUserByEmail with valid email
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test GetUserByEmail with valid email
	user, err := userService.GetUserByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Errorf("GetUserByEmail failed: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("GetUserByEmail returned incorrect email: %s", user.Email)
	}
}

func TestUserService_GetUserByEmail_UserNotFound(t *testing.T) {
	// Test GetUserByEmail with user not found
	repo := &mockUserRepository{
		err: pkg.ErrUserNotFound,
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test GetUserByEmail with user not found
	_, err := userService.GetUserByEmail(context.Background(), "test@example.com")
	fmt.Println(err)
	if err != pkg.ErrUserNotFound {
		t.Errorf("GetUserByEmail should have failed with user not found: %v", err)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	// Test GetAllUsers with valid users
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test GetAllUsers with valid users
	users, err := userService.GetAllUsers(context.Background())
	if err != nil {
		t.Errorf("GetAllUsers failed: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("GetAllUsers returned incorrect number of users: %d", len(users))
	}
}

func TestUserService_GetAllUsers_Error(t *testing.T) {
	// Test GetAllUsers with error
	repo := &mockUserRepository{
		err: errors.New("test error"),
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test GetAllUsers with error
	_, err := userService.GetAllUsers(context.Background())
	fmt.Println(err)
	if err == nil || err.Error() != "test error" {
		t.Errorf("GetAllUsers should have failed with error: %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	// Test DeleteUser with valid user ID
	repo := &mockUserRepository{
		users: []domain.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test DeleteUser with valid user ID
	err := userService.DeleteUser(context.Background(), 1)
	if err != nil {
		t.Errorf("DeleteUser failed: %v", err)
	}
}

func TestUserService_DeleteUser_UserNotFound(t *testing.T) {
	// Test DeleteUser with user not found
	repo := &mockUserRepository{
		err: pkg.ErrUserNotFound,
	}
	tokenProvider := &mockTokenProvider{
		token: "token",
		sub:   "sub",
	}
	userService := NewUserService(repo, tokenProvider)
	// Test DeleteUser with user not found
	err := userService.DeleteUser(context.Background(), 2)
	fmt.Println(err)
	if err != pkg.ErrUserNotFound {
		t.Errorf("DeleteUser should have failed with user not found: %v", err)
	}
}
