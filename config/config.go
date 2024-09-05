// package config

// import (
//     "fmt"
//     "os"
//     "gorm.io/driver/postgres"
//     "gorm.io/gorm"
// //    "log"
// )

// var DB *gorm.DB

// func InitDB() {
//     dsn := fmt.Sprintf(
//         "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
//         os.Getenv("DB_HOST"),
//         os.Getenv("DB_USER"),
//         os.Getenv("DB_PASSWORD"),
//         os.Getenv("DB_NAME"),
//         os.Getenv("DB_PORT"),
//     )
//     var err error
//     DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//     if err != nil {
//         panic("failed to connect to database")
//     }
// }

package config

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return DB, nil
}
