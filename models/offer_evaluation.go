package models

import "time"

type OfferEvaluation struct {
    ID          uint    `gorm:"primaryKey"`
    OfferID     uint    `json:"offer_id"`  // Referência para a oferta
    Offer       Offer   `gorm:"foreignKey:OfferID"`
    UserID      uint    `json:"user_id"`   // Referência para o cliente que avaliou
    User        User    `gorm:"foreignKey:UserID"`
    CouponID    uint    `json:"coupon_id"` // Referência para o cupom resgatado
    Coupon      Coupon  `gorm:"foreignKey:CouponID"`
    Rating      float64 `json:"rating" gorm:"default:0.0"`    // Nota dada pelo cliente (exemplo: 1 a 5 estrelas)
    Comment     string  `json:"comment"`   // Comentário opcional do cliente
    CreatedAt   time.Time `json:"created_at"`
}
