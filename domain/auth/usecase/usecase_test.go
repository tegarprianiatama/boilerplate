package usecase

import (
	"boilerplate/model/entity"
	"boilerplate/model/request"
	"boilerplate/shared/auth"
	ce "boilerplate/shared/common/errors"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUsecase_Register(t *testing.T) {
	mockRepo := new(MockRepository)
	u := NewUsecase(mockRepo)

	ctx := context.Background()
	req := request.CreateUserReq{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
		Role:     "user",
	}

	mockRepo.On("FindByEmail", ctx, req.Email).Return(entity.User{}, gorm.ErrRecordNotFound)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(user *entity.User) bool {
		return user.Password != "" && user.Password != req.Password
	})).Return(nil)

	res, err := u.Register(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "John Doe", res.Name)
	assert.Equal(t, "john@example.com", res.Email)
	mockRepo.AssertExpectations(t)
}

func TestUsecase_Login(t *testing.T) {
	mockRepo := new(MockRepository)
	u := NewUsecase(mockRepo)

	ctx := context.Background()
	hashedPassword, err := auth.HashPassword("secret123")
	assert.NoError(t, err)

	existingUser := entity.User{
		ID:       1,
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: hashedPassword,
		Role:     "admin",
	}

	mockRepo.On("FindByEmail", ctx, existingUser.Email).Return(existingUser, nil)

	res, err := u.Login(ctx, request.LoginReq{
		Email:    existingUser.Email,
		Password: "secret123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "Bearer", res.TokenType)
	assert.Equal(t, existingUser.ID, res.User.ID)
	assert.Equal(t, existingUser.Email, res.User.Email)
	assert.False(t, res.ExpiresAt.IsZero())
	mockRepo.AssertExpectations(t)
}

func TestUsecase_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockRepository)
	u := NewUsecase(mockRepo)

	ctx := context.Background()
	hashedPassword, err := auth.HashPassword("secret123")
	assert.NoError(t, err)

	existingUser := entity.User{
		ID:       1,
		Email:    "jane@example.com",
		Password: hashedPassword,
		Role:     "admin",
	}

	mockRepo.On("FindByEmail", ctx, existingUser.Email).Return(existingUser, nil)

	_, err = u.Login(ctx, request.LoginReq{
		Email:    existingUser.Email,
		Password: "wrong-password",
	})

	assert.ErrorIs(t, err, ce.ErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
}
