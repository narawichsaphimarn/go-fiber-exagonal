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
	// Validate user input
	if err := pkg.ValidateStruct(ctx, user); err != nil {
		return err
	}
	// Check if user already exists
	userDB, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if userDB != nil && userDB.ID > 0 {
		return pkg.ErrUserAlreadyExists
	}
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
	if userDB == nil || userDB.ID == 0 {
		return "", pkg.ErrUserNotFound
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

func (s *UserService) UpdateUser(ctx context.Context, id int, user *domain.UpdateUserRequest) error {
	// Validate user input
	if err := pkg.ValidateStruct(ctx, user); err != nil {
		return err
	}
	userDB, err := s.repo.GetUserByID(ctx, strconv.Itoa(id))
	if err != nil {
		return err
	}
	if userDB == nil || userDB.ID == 0 {
		return pkg.ErrUserNotFound
	}
	return s.repo.UpdateUser(ctx, &domain.User{
		ID:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) UpdatePassword(ctx context.Context, id int, req domain.UpdatePasswordRequest) error {
	// Validate password
	if err := pkg.ValidateStruct(ctx, req); err != nil {
		return err
	}
	// Check if user exists
	userDB, err := s.repo.GetUserByID(ctx, strconv.Itoa(id))
	if err != nil {
		return err
	}
	if userDB == nil || userDB.ID == 0 {
		return pkg.ErrUserNotFound
	}
	// Hash password
	hashedPassword, err := pkg.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, id, hashedPassword)
}
