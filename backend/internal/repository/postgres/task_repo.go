package postgres

import (
	"context"

	"taskmaster/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, priority, assignee_id, creator_id, team_id, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query, task.ID, task.Title, task.Description, task.Status, task.Priority, task.AssigneeID, task.CreatorID, task.TeamID, task.DueDate, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *TaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	query := `SELECT id, title, description, status, priority, assignee_id, creator_id, team_id, due_date, created_at, updated_at FROM tasks WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var t domain.Task
	var assigneeID, dueDate *interface{}
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &assigneeID, &t.CreatorID, &t.TeamID, &dueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) FindByTeam(ctx context.Context, teamID uuid.UUID, limit, offset int) ([]domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, assignee_id, creator_id, team_id, due_date, created_at, updated_at
		FROM tasks WHERE team_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, teamID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		var assigneeID, dueDate *interface{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &assigneeID, &t.CreatorID, &t.TeamID, &dueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) CountByTeam(ctx context.Context, teamID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM tasks WHERE team_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, teamID).Scan(&count)
	return count, err
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, status=$3, priority=$4, assignee_id=$5, due_date=$6, updated_at=$7 WHERE id=$8`
	_, err := r.db.Exec(ctx, query, task.Title, task.Description, task.Status, task.Priority, task.AssigneeID, task.DueDate, task.UpdatedAt, task.ID)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	return err
}

func (r *TaskRepository) Assign(ctx context.Context, taskID, assigneeID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE tasks SET assignee_id = $1, updated_at = NOW() WHERE id = $2`, assigneeID, taskID)
	return err
}
