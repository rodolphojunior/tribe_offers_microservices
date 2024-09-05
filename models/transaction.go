package models

import (
    "time"
)

type TransactionStatus string

const (
    StatusPending   TransactionStatus = "pending"
    StatusSuccess   TransactionStatus = "success"
    StatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
    ID            uint              `json:"id" gorm:"primaryKey"`
    UserID        uint              `json:"user_id"`
    User          User              `json:"user"`
    OfferID       uint              `json:"offer_id"`
    Offer         Offer             `json:"offer"`
    Amount        float64           `json:"amount"`
    Status        TransactionStatus `json:"status"` // Ex: pending, completed, failed
    CreatedAt     time.Time         `json:"created_at"`
    UpdatedAt     time.Time         `json:"updated_at"`
}
