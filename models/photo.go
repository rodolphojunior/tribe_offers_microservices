package models

import (
    "time"
)

type Photo struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    OfferID   uint      `json:"offer_id"`            // Chave estrangeira para Offer
    Offer     Offer     `json:"offer"`               // Referência ao modelo Offer
    FilePath  string    `json:"file_path"`           // URL da imagem
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
