package postgres

import (
	"context"

	"taskmaster/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password, name, role, team_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Password, user.Name, user.Role, user.TeamID, user.CreatedAt)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password, name, role, team_id, created_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user domain.User
	var teamID *uuid.UUID
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &teamID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	user.TeamID = teamID
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, password, name, role, team_id, created_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user domain.User
	var teamID *uuid.UUID
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Role, &teamID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	user.TeamID = teamID
	return &user, nil
}
