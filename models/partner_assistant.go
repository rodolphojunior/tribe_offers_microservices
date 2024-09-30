package models

import (
    "time"
)

type PartnerAssistant struct {
    ID           uint         `json:"id" gorm:"primaryKey"`
    UserID       uint         `json:"user_id"`
    User         User         `json:"user" gorm:"foreignKey:UserID"`
    PartnerAdmin PartnerAdmin `json:"partner_admin" gorm:"foreignKey:PartnerAdminID"`
    PartnerAdminID uint       `json:"partner_admin_id"`  // Campo de chave estrangeira
    CompanyID    uint         `json:"company_id"`
    Company      Company      `json:"company" gorm:"foreignKey:CompanyID"`
    CreatedAt    time.Time    `json:"created_at"`
    UpdatedAt    time.Time    `json:"updated_at"`
}
