package models

import (
    "time"
)
gi
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username"`
    FullName  string    `json:"full_name"`
    Email     string    `json:"email"`
    Password  string    `json:"password"`
    Role      string    `json:"role"`
    Enabled   bool      `json:"enabled" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Cpf       string    `json:"cpf"`
}
