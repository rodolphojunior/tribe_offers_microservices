package models

import (
    "time"
)

type Company struct {
    ID              uint        `json:"id" gorm:"primaryKey"`
    Cnpj            string      `json:"cnpj" gorm:"unique;not null"`
    CompanyName     string      `json:"company_name"`
    Description     string      `json:"description"`
    TradeName       string      `json:"trade_name"` 
    Cep             string      `json:"cep"`
    Address         string      `json:"address"`
    City            string      `json:"city"`
    State           string      `json:"state"`
    PhoneNumber     string      `json:"phone_number"`
    Email           string      `json:"email"`
    Website         string      `json:"website"` 
    BankName        string      `json:"bank_name"`  
    AgencyWithDigit string      `json:"agency_with_digit"`   
    CurrentAccount  string      `json:"current_account"`  
    DigitAccount    string      `json:"digit_account"`  
    CreatedAt       time.Time   `json:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at"`
    Offers          []Offer     `json:"offers" gorm:"foreignKey:CompanyID"`
}
