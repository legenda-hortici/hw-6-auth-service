package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	ID        int       `gorm:"primaryKey;column:id" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id;not null" json:"user_id"`
	Hash      uuid.UUID `gorm:"column:token_hash;not null;uniqueIndex" json:"hash"`
	ExpiresAt int64     `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	Hash      uuid.UUID `json:"hash"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ExpireAt  int64     `json:"expire_at"`
}
