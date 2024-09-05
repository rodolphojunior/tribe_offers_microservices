package models

type Partner struct {
    ID        uint   `json:"id" gorm:"primaryKey"`
    UserID    uint   `json:"user_id"`
    User      User   `json:"user" gorm:"foreignKey:UserID"`
    CompanyID uint   `json:"company_id"`
    Company   Company `json:"company" gorm:"foreignKey:CompanyID"`
}