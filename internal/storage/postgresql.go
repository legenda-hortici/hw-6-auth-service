package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/legenda-hortici/hw-6-auth-service/internal/config"
	"github.com/legenda-hortici/hw-6-auth-service/internal/domain"
	"github.com/legenda-hortici/hw-6-auth-service/internal/storage/myerr"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(cfg config.Config) (*Storage, error) {
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
		panic("Failed to connect to database")
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CheckUser(ctx context.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user domain.Users
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

	user := domain.Users{
		ID:       uuid.New(),
		Username: username,
		Password: passHash,
	}

	result := s.db.Create(&user)
	if result.Error != nil {
		return myerr.FailedToCreateUserErr
	}

	return nil
}

func (s *Storage) Login(ctx context.Context, email string) (*domain.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user domain.Users

	err := s.db.Where("username = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerr.UserNotFoundErr
		}
		return nil, errors.Wrap(err, "failed to query user")
	}

	return &user, nil
}
