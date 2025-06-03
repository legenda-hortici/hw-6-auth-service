package services

import (
	"context"
	"github.com/pkg/errors"
	"testing"

	"github.com/legenda-hortici/hw-6-auth-service/internal/config"
	"github.com/legenda-hortici/hw-6-auth-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockAuthRepository реализует AuthRepository для тестов
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Register(ctx context.Context, username string, passHash []byte) error {
	args := m.Called(ctx, username, passHash)
	return args.Error(0)
}

func (m *MockAuthRepository) Login(ctx context.Context, email string) (*domain.Users, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.Users), args.Error(1)
}

func (m *MockAuthRepository) CheckUser(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		checkUser     func(ctx context.Context, username string) (bool, error)
		register      func(ctx context.Context, username string, passHash []byte) error
		expectedError error
	}{
		{
			name:     "successful registration",
			username: "testuser",
			password: "testpass",
			checkUser: func(ctx context.Context, username string) (bool, error) {
				return false, nil
			},
			register: func(ctx context.Context, username string, passHash []byte) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name:     "user already exists",
			username: "existinguser",
			password: "testpass",
			checkUser: func(ctx context.Context, username string) (bool, error) {
				return true, nil
			},
			register:      nil, // не должно вызываться
			expectedError: errors.New("User already exists"),
		},
		{
			name:     "check user error",
			username: "testuser",
			password: "testpass",
			checkUser: func(ctx context.Context, username string) (bool, error) {
				return false, errors.New("database error")
			},
			register:      nil, // не должно вызываться
			expectedError: errors.New("database error"),
		},
		{
			name:     "register error",
			username: "testuser",
			password: "testpass",
			checkUser: func(ctx context.Context, username string) (bool, error) {
				return false, nil
			},
			register: func(ctx context.Context, username string, passHash []byte) error {
				return errors.New("registration failed")
			},
			expectedError: errors.New("registration failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем mock репозитория
			mockRepo := new(MockAuthRepository)

			// Настраиваем ожидания для CheckUser
			mockRepo.On("CheckUser", mock.Anything, tt.username).
				Return(tt.checkUser(context.Background(), tt.username))

			// Если register не nil, настраиваем ожидания
			if tt.register != nil {
				// Для проверки хеша пароля
				mockRepo.On("Register", mock.Anything, tt.username, mock.AnythingOfType("[]uint8")).
					Return(tt.register(context.Background(), tt.username, []byte{}))
			}

			// Создаем сервис с mock репозиторием
			service := NewAuthService(
				config.Config{},
				zap.NewNop().Sugar(),
				mockRepo,
			)

			// Вызываем тестируемый метод
			err := service.Register(context.Background(), tt.username, tt.password)

			// Проверяем ошибку
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем, что все ожидания выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}
