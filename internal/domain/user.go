package domain

import "github.com/google/uuid"

type Users struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password []byte    `json:"password"`
}
