package usercase

import (
	"context"
	"strconv"

	"example.com/practice/fiber/internal/domain"
	"example.com/practice/fiber/internal/ports"
	"example.com/practice/fiber/pkg"
)

type UserService struct {
	repo          ports.UserRepository
	tokenProvider ports.TokenProvider
}

func NewUserService(repo ports.UserRepository, tokenProvider ports.TokenProvider) *UserService {
	return &UserService{repo: repo, tokenProvider: tokenProvider}
}

func (s *UserService) Register(ctx context.Context, user *domain.User) error {
	// Hash password
	hashedPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.Role = "user"
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) Login(ctx context.Context, user *domain.AuthRequest) (string, error) {
	// Validate user input
	if err := pkg.ValidateStruct(ctx, user); err != nil {
		return "", err
	}
	// Get user by email
	userDB, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	// Compare password
	err = userDB.ComparePassword(user.Password)
	if err != nil {
		return "", err
	}
	// Generate token
	token, err := s.tokenProvider.GenerateToken(strconv.Itoa(userDB.ID))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) UpdatePassword(ctx context.Context, id int, password string) error {
	// Hash password
	hashedPassword, err := pkg.HashPassword(password)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, id, hashedPassword)
}
