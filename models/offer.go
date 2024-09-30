package models

import (
    "time"
)

type Offer struct {
	ID            	uint      `json:"id" gorm:"primaryKey"`
	Price         	float64   `json:"price"`
    Title  		  	string    `json:"title"`
    Description   	string    `json:"description"`
	Discount      	float64   `json:"discount"`
	Commission    	float64   `json:"commission"`
	Tariff	      	float64   `json:"tariff"`	
	StartDate     	time.Time `json:"start_date"`
	EndDate       	time.Time `json:"end_date"`
	PromoUnits    	uint64    `json:"promo_units"`
	UnitsSold     	uint64    `json:"units_sold"`
	AvailableUnits 	uint64   `json:"available_units"`
	Enable        	bool      `json:"enable"`
	R1            	float64   `json:"r1"`     // Calculado: Discount /  Price
	R2            	float64   `json:"r2"`     // Calculado: Commission / Price
	R3            	float64   `json:"r3"`     // Calculado: UnitsSold / PromoUnits
	R4            	float64   `json:"r4"`     // Calculado: Dias Restantes / Dias Totais
	R5            	float64   `json:"r5"`     // Calculado: r5 = average_rating / 5
	YNorm         	float64   `json:"y_norm"` // Calculado: Normalizado para ordenação
	AverageRating 	float64   `json:"average_rating" gorm:"default:0.0"` // Valor padrão de 0.0
	CreatedAt     	time.Time `json:"created_at"`
	UpdatedAt     	time.Time `json:"updated_at"`
    CompanyID     	uint      `json:"company_id"`
    Company       	Company   `json:"company"`
    Photos        	[]Photo   `json:"photos" gorm:"foreignKey:OfferID"` // Relação um-para-muitos com o modelo Photo    
}