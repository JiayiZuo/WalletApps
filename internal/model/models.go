package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	PasswordHash string
	Salt         string
	CreatedAt    time.Time
}

type Wallet struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Address   string
	Balance   float64
	UpdatedAt time.Time
}

type Transaction struct {
	ID            uuid.UUID
	WalletID      uuid.UUID
	Amount        float64
	Type          string
	Description   string
	RelatedUserID *uuid.UUID
	CreatedAt     time.Time
}
