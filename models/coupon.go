package models

import (
    "errors"
    "gorm.io/gorm"
    "time"
)

type CouponStatus string

const (
    StatusActive  CouponStatus = "active"
    StatusExpired CouponStatus = "expired"
    StatusUsed    CouponStatus = "used"
)

type Coupon struct {
    ID          uint         `json:"id" gorm:"primaryKey"`
    Code        string       `json:"code"`           // Código do cupom gerado
    Status      CouponStatus `json:"status"`         // Ex: active, expired ou used???
    Paid        bool         `json:"paid"`           // Confirmado como pago
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    OfferID     uint         `json:"offer_id"`
    Offer       Offer        `json:"offer"`
    UserID      uint         `json:"user_id"`
    User        User         `json:"user"`
    // Removed TransactionID - não é necessário
}

// BeforeSave é chamado antes de salvar o registro no banco de dados
func (c *Coupon) BeforeSave(tx *gorm.DB) (err error) {
    if c.Status != StatusActive && c.Status != StatusExpired && c.Status != StatusUsed {
        return errors.New("invalid status")
    }
    return
}
