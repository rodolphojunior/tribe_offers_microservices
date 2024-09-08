package models

import (
    "time"
)

type Offer struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Price         float64   `json:"price"`
	Discount      float64   `json:"discount"`
	Commission    float64   `json:"commission"`
	Tariff	      float64   `json:"tariff"`	
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	PromoUnits    uint64    `json:"promo_units"`
	UnitsSold     uint64    `json:"units_sold"`
	AvaialbleUnits uint64   `json:"available_units"`
	Enable        bool      `json:"enable"`
	R1            float64   `json:"r1"`     // Calculado: Discount /  Price
	R2            float64   `json:"r2"`     // Calculado: Commission / Price
	R3            float64   `json:"r3"`     // Calculado: UnitsSold / PromoUnits
	R4            float64   `json:"r4"`     // Calculado: Dias Restantes / Dias Totais
	YNorm         float64   `json:"y_norm"` // Calculado: Normalizado para ordenação
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
    CompanyID     uint      `json:"company_id"`
    Company       Company   `json:"company"`
    Coupons       []Coupon  `json:"coupons" gorm:"foreignKey:OfferID"` // Adicionando o campo Coupons
    Photos        []Photo   `json:"photos" gorm:"foreignKey:OfferID"` // Relação um-para-muitos com o modelo Photo    
}