package unit

import (
	"context"
	"testing"
	"time"

	"taskmaster/internal/domain"
	"taskmaster/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByTeam(ctx context.Context, teamID uuid.UUID, limit, offset int) ([]domain.Task, error) {
	args := m.Called(ctx, teamID, limit, offset)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) CountByTeam(ctx context.Context, teamID uuid.UUID) (int, error) {
	args := m.Called(ctx, teamID)
	return args.Int(0), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) Assign(ctx context.Context, taskID, assigneeID uuid.UUID) error {
	args := m.Called(ctx, taskID, assigneeID)
	return args.Error(0)
}

// MockTeamRepository
type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) IsMember(ctx context.Context, teamID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, teamID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTeamRepository) GetSLA(ctx context.Context, teamID uuid.UUID) (int, error) {
	args := m.Called(ctx, teamID)
	return args.Int(0), args.Error(1)
}

func (m *MockTeamRepository) GetTaskLimit(ctx context.Context, teamID uuid.UUID) (int, error) {
	args := m.Called(ctx, teamID)
	return args.Int(0), args.Error(1)
}

// MockWebSocketHub
type MockWebSocketHub struct {
	mock.Mock
}

func (m *MockWebSocketHub) Broadcast(teamID uuid.UUID, message interface{}) {
	m.Called(teamID, message)
}

func (m *MockWebSocketHub) Register(client *domain.User, conn interface{}) {
	m.Called(client, conn)
}

func (m *MockWebSocketHub) Unregister(client *domain.User) {
	m.Called(client)
}

func TestTaskService_Create_Success(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	teamRepo := new(MockTeamRepository)
	cacheRepo := new(MockCacheRepository)
	wsHub := new(MockWebSocketHub)

	svc := service.NewTaskService(taskRepo, teamRepo, cacheRepo, wsHub)

	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	creatorID := uuid.New()

	teamRepo.On("GetTaskLimit", mock.Anything, teamID).Return(100, nil)
	taskRepo.On("CountByTeam", mock.Anything, teamID).Return(5, nil)
	teamRepo.On("GetSLA", mock.Anything, teamID).Return(14, nil)
	taskRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil)
	cacheRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)
	wsHub.On("Broadcast", teamID, mock.Anything).Return()

	dueDate := time.Now().AddDate(0, 0, 7)
	task, err := svc.Create(context.Background(), service.CreateTaskInput{
		Title:    "Test Task",
		Priority: domain.PriorityHigh,
		DueDate:  &dueDate,
	}, creatorID, teamID)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, domain.TaskStatusTodo, task.Status)
	assert.Equal(t, creatorID, task.CreatorID)
	assert.Equal(t, teamID, task.TeamID)

	taskRepo.AssertExpectations(t)
	teamRepo.AssertExpectations(t)
	wsHub.AssertExpectations(t)
}

func TestTaskService_Create_ExceedsLimit(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	teamRepo := new(MockTeamRepository)
	cacheRepo := new(MockCacheRepository)
	wsHub := new(MockWebSocketHub)

	svc := service.NewTaskService(taskRepo, teamRepo, cacheRepo, wsHub)

	teamID := uuid.New()
	creatorID := uuid.New()

	teamRepo.On("GetTaskLimit", mock.Anything, teamID).Return(10, nil)
	taskRepo.On("CountByTeam", mock.Anything, teamID).Return(10, nil)

	_, err := svc.Create(context.Background(), service.CreateTaskInput{
		Title:    "Overflow Task",
		Priority: domain.PriorityMedium,
	}, creatorID, teamID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit reached")
}

func TestTaskService_Create_InvalidAssignee(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	teamRepo := new(MockTeamRepository)
	cacheRepo := new(MockCacheRepository)
	wsHub := new(MockWebSocketHub)

	svc := service.NewTaskService(taskRepo, teamRepo, cacheRepo, wsHub)

	teamID := uuid.New()
	creatorID := uuid.New()
	assigneeID := uuid.New()

	teamRepo.On("IsMember", mock.Anything, teamID, assigneeID).Return(false, nil)

	_, err := svc.Create(context.Background(), service.CreateTaskInput{
		Title:      "Bad Assignee",
		Priority:   domain.PriorityLow,
		AssigneeID: &assigneeID,
	}, creatorID, teamID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a team member")
}
