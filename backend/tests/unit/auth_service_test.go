package unit

import (
	"context"
	"testing"

	"taskmaster/internal/domain"
	"taskmaster/internal/ports"
	"taskmaster/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockCacheRepository mocks the CacheRepository interface
type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	args := m.Called(ctx, key, value, ttlSeconds)
	return args.Error(0)
}

func (m *MockCacheRepository) Delete(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

func (m *MockCacheRepository) Increment(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheRepository) Expire(ctx context.Context, key string, seconds int) error {
	args := m.Called(ctx, key, seconds)
	return args.Error(0)
}

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	cacheRepo := new(MockCacheRepository)
	svc := service.NewAuthService(userRepo, cacheRepo, "test-secret")

	userRepo.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, assert.AnError)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := svc.Register(context.Background(), service.RegisterInput{
		Email:    "new@example.com",
		Password: "password123",
		Name:     "Test User",
	})

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, domain.RoleMember, user.Role)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	userRepo := new(MockUserRepository)
	cacheRepo := new(MockCacheRepository)
	svc := service.NewAuthService(userRepo, cacheRepo, "test-secret")

	existing := &domain.User{ID: uuid.New(), Email: "exists@example.com"}
	userRepo.On("FindByEmail", mock.Anything, "exists@example.com").Return(existing, nil)

	_, err := svc.Register(context.Background(), service.RegisterInput{
		Email:    "exists@example.com",
		Password: "password123",
		Name:     "Test",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	cacheRepo := new(MockCacheRepository)
	svc := service.NewAuthService(userRepo, cacheRepo, "test-secret")

	userRepo.On("FindByEmail", mock.Anything, "notfound@example.com").Return(nil, assert.AnError)

	_, _, err := svc.Login(context.Background(), service.LoginInput{
		Email:    "notfound@example.com",
		Password: "wrong",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}
