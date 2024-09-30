package models

import (
    "time"
)

// type Partner struct {
//     ID        uint   `json:"id" gorm:"primaryKey"`
//     UserID    uint   `json:"user_id"`
//     User      User   `json:"user" gorm:"foreignKey:UserID"`
//     CompanyID uint   `json:"company_id"`
//     Company   Company `json:"company" gorm:"foreignKey:CompanyID"`
// }

type PartnerAdmin struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    UserID      uint      `json:"user_id"`
    User        User      `json:"user" gorm:"foreignKey:UserID"`
    CompanyID   uint      `json:"company_id" gorm:"unique"`
    Company     Company   `json:"company" gorm:"foreignKey:CompanyID"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}