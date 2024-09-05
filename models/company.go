package models

import (
    "time"
)

type Company struct {
    ID              uint        `json:"id" gorm:"primaryKey"`
    Cnpj            string      `json:"cnpj" gorm:"unique;not null"`
    TradeName       string      `json:"trade_name"`  // Novo campo
    Description     string      `json:"description"`
    Address         string      `json:"address"`
    PhoneNumber     string      `json:"phone_number"`
    Email           string      `json:"email"`
    Website         string      `json:"website"` 
    CompanyName     string      `json:"company_name"`  // Novo campo
    BankName        string      `json:"bank_name"`  // Novo campo
    AgencyWithDigit string      `json:"agency_with_digit"`   // Novo campo
    CurrentAccount  string      `json:"current_account"`  // Novo campo
    DigitAccount    string      `json:"digit_account"`  // Novo campo
    CreatedAt       time.Time   `json:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at"`
    Offers          []Offer     `json:"offers" gorm:"foreignKey:CompanyID"`
}
