package ports

import (
	"context"

	"taskmaster/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	FindByTeam(ctx context.Context, teamID uuid.UUID, limit, offset int) ([]domain.Task, error)
	CountByTeam(ctx context.Context, teamID uuid.UUID) (int, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	Assign(ctx context.Context, taskID, assigneeID uuid.UUID) error
}

type TeamRepository interface {
	IsMember(ctx context.Context, teamID, userID uuid.UUID) (bool, error)
	GetSLA(ctx context.Context, teamID uuid.UUID) (int, error)
	GetTaskLimit(ctx context.Context, teamID uuid.UUID) (int, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttlSeconds int) error
	Delete(ctx context.Context, pattern string) error
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, seconds int) error
}

type WebSocketHub interface {
	Broadcast(teamID uuid.UUID, message interface{})
	Register(client *domain.User, conn interface{})
	Unregister(client *domain.User)
}
