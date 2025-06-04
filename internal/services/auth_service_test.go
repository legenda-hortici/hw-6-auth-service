package services_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/legenda-hortici/hw-6-auth-service/internal/config"
	"github.com/legenda-hortici/hw-6-auth-service/internal/domain"
	"github.com/legenda-hortici/hw-6-auth-service/internal/mocks"
	"github.com/legenda-hortici/hw-6-auth-service/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

//mockery --name=AuthRepository --dir=internal/services --output=internal/mocks --log-level=debug

func TestAuthService_Register(t *testing.T) {
	// создаем слой репозитория
	mockRepo := mocks.NewAuthRepository(t)

	// настраиваем поведение метода CheckUser (то есть пользователь не существует)
	mockRepo.On("CheckUser", mock.Anything, "user1").Return(false, nil)

	// поведение для метода Register - успешная регистрация
	mockRepo.On("Register", mock.Anything, "user1", mock.Anything).Return(nil)

	svc := services.NewAuthService(config.Config{}, nil, mockRepo)

	err := svc.Register(context.Background(), "user1", "password")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := mocks.NewAuthRepository(t)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mockRepo.On("Login", mock.Anything, "user1").Return(&domain.Users{
		ID:       uuid.New(),
		Username: "user1",
		Password: hashed,
	}, nil)

	logger := zap.NewExample().Sugar()

	svc := services.NewAuthService(config.Config{}, logger, mockRepo)

	token, err := svc.Login(context.Background(), "user1", "password")

	assert.NoError(t, err)

	assert.NotEmpty(t, token)
}
