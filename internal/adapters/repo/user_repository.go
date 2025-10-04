package repo

import (
	"context"

	domain "example.com/practice/fiber/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetUserByID(ctx context.Context, userId string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, "SELECT id, email, password, username, first_name, last_name, role FROM users WHERE id = $1", userId).Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.FirstName, &user.LastName, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {	
	var user domain.User
	err := r.db.QueryRow(ctx, "SELECT id, email, password, username, first_name, last_name, role FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.FirstName, &user.LastName, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	rows, err := r.db.Query(ctx, "SELECT id, email, password, username, first_name, last_name, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.FirstName, &user.LastName, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *domain.User) error {		
	_, err := r.db.Exec(ctx, "INSERT INTO users (email, password, username, first_name, last_name, role) VALUES ($1, $2, $3, $4, $5, $6)", user.Email, user.Password, user.Username, user.FirstName, user.LastName, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.Exec(ctx, "UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3", user.FirstName, user.LastName, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id int) error {						
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, id int, password string) error {
	_, err := r.db.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", password, id)
	if err != nil {
		return err
	}
	return nil
}
