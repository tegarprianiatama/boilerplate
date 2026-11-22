package usecase

import (
	"boilerplate/model/entity"
	"boilerplate/model/request"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of interfaces.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id uint) (entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUsecase_Create(t *testing.T) {
	mockRepo := new(MockRepository)
	u := NewUsecase(mockRepo)

	ctx := context.Background()
	req := request.CreateUserReq{
		Name:  "John Doe",
		Email: "john@example.com",
		Role:  "user",
	}

	// Mock FindByEmail returns gorm.ErrRecordNotFound (email available)
	mockRepo.On("FindByEmail", ctx, req.Email).Return(entity.User{}, gorm.ErrRecordNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	res, err := u.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "John Doe", res.Name)
	assert.Equal(t, "john@example.com", res.Email)
	mockRepo.AssertExpectations(t)
}

func TestUsecase_GetByID(t *testing.T) {
	mockRepo := new(MockRepository)
	u := NewUsecase(mockRepo)

	ctx := context.Background()
	expectedUser := entity.User{
		ID:        1,
		Name:      "Jane Doe",
		Email:     "jane@example.com",
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(expectedUser, nil)

	res, err := u.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), res.ID)
	assert.Equal(t, "Jane Doe", res.Name)
	mockRepo.AssertExpectations(t)
}
