package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"taskmaster/internal/domain"
	"taskmaster/internal/ports"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo  ports.TaskRepository
	teamRepo  ports.TeamRepository
	cacheRepo ports.CacheRepository
	wsHub     ports.WebSocketHub
}

func NewTaskService(taskRepo ports.TaskRepository, teamRepo ports.TeamRepository, cacheRepo ports.CacheRepository, wsHub ports.WebSocketHub) *TaskService {
	return &TaskService{taskRepo, teamRepo, cacheRepo, wsHub}
}

type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    domain.Priority `json:"priority" binding:"required"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

func (s *TaskService) Create(ctx context.Context, input CreateTaskInput, creatorID, teamID uuid.UUID) (*domain.Task, error) {
	// Check team membership of assignee
	if input.AssigneeID != nil {
		isMember, err := s.teamRepo.IsMember(ctx, teamID, *input.AssigneeID)
		if err != nil || !isMember {
			return nil, errors.New("assignee is not a team member")
		}
	}

	// Check task limit
	limit, err := s.teamRepo.GetTaskLimit(ctx, teamID)
	if err != nil {
		return nil, err
	}
	count, err := s.taskRepo.CountByTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if count >= limit {
		return nil, fmt.Errorf("team task limit reached (%d/%d)", count, limit)
	}

	// Check SLA
	if input.DueDate != nil {
		sla, err := s.teamRepo.GetSLA(ctx, teamID)
		if err != nil {
			return nil, err
		}
		maxDue := time.Now().AddDate(0, 0, sla)
		if input.DueDate.After(maxDue) {
			return nil, fmt.Errorf("due date exceeds team SLA of %d days", sla)
		}
	}

	task := &domain.Task{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		Status:      domain.TaskStatusTodo,
		Priority:    input.Priority,
		AssigneeID:  input.AssigneeID,
		CreatorID:   creatorID,
		TeamID:      teamID,
		DueDate:     input.DueDate,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Invalidate cache
	_ = s.cacheRepo.Delete(ctx, fmt.Sprintf("tasks:team:%s:*", teamID.String()))

	// Real-time broadcast
	s.wsHub.Broadcast(teamID, map[string]interface{}{
		"type":    "task:created",
		"payload": task,
	})

	return task, nil
}

func (s *TaskService) List(ctx context.Context, teamID uuid.UUID, limit, offset int) ([]domain.Task, error) {
	return s.taskRepo.FindByTeam(ctx, teamID, limit, offset)
}

func (s *TaskService) Get(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return s.taskRepo.FindByID(ctx, id)
}

func (s *TaskService) Delete(ctx context.Context, id, teamID uuid.UUID) error {
	if err := s.taskRepo.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cacheRepo.Delete(ctx, fmt.Sprintf("tasks:team:%s:*", teamID.String()))
	s.wsHub.Broadcast(teamID, map[string]interface{}{
		"type": "task:deleted",
		"payload": map[string]string{"id": id.String()},
	})
	return nil
}
