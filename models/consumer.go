package models

type Consumer struct {
    ID     uint `json:"id" gorm:"primaryKey"`
    UserID uint `json:"user_id"`
    User   User `json:"user" gorm:"foreignKey:UserID"`
    Coupons   []Coupon `json:"coupons" gorm:"foreignKey:UserID"` // Relacionamento com cupons
}