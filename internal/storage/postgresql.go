package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"skillsRockAuthService/internal/config"
	"skillsRockAuthService/internal/domain"
	"skillsRockAuthService/internal/storage/myerr"
	"time"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(log *zap.SugaredLogger, cfg config.Config) (*Storage, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DatabaseName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CheckUser(ctx context.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user string

	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, errors.Wrap(err, "failed to check user")
	}

	return true, myerr.AlreadyExistsErr
}

func (s *Storage) Register(ctx context.Context, username string, passHash []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := domain.User{
		ID:       uuid.New(),
		Username: username,
		PassHash: passHash,
	}

	result := s.db.Create(&user)
	if result.Error != nil {
		return myerr.FailedToCreateUserErr
	}

	return nil
}

func (s *Storage) Login(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user domain.User

	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerr.UserNotFoundErr
		}
		return nil, errors.Wrap(err, "failed to query user")
	}

	return &user, nil
}
