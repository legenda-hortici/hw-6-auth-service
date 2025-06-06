package domain

import "github.com/google/uuid"

type Users struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username string    `json:"username"`
	Password []byte    `json:"password"`
}
