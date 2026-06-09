package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db}
}

func (r *TeamRepository) IsMember(ctx context.Context, teamID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, teamID, userID).Scan(&exists)
	return exists, err
}

func (r *TeamRepository) GetSLA(ctx context.Context, teamID uuid.UUID) (int, error) {
	query := `SELECT COALESCE(sla_days, 14) FROM teams WHERE id = $1`
	var sla int
	err := r.db.QueryRow(ctx, query, teamID).Scan(&sla)
	return sla, err
}

func (r *TeamRepository) GetTaskLimit(ctx context.Context, teamID uuid.UUID) (int, error) {
	query := `SELECT COALESCE(task_limit, 50) FROM teams WHERE id = $1`
	var limit int
	err := r.db.QueryRow(ctx, query, teamID).Scan(&limit)
	return limit, err
}
