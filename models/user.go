package models

import (
    "time"
)

type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"password"`
    FullName  string    `json:"full_name"`
    Role      string    `json:"role"`
    Enabled   bool      `json:"enabled" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
