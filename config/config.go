package config

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "log"
    "fmt"
    "os"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
    // Verifique as variáveis de ambiente e imprima mensagens de erro se estiverem faltando
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
        log.Fatal("Uma ou mais variáveis de ambiente do banco de dados não foram definidas")
    }

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        dbHost, dbUser, dbPassword, dbName, dbPort,
    )

    log.Println("Conectando ao banco de dados com DSN:", dsn)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
        return nil, err
    }

    DB = db.Debug() // Habilitar Debug para ver as consultas SQL
    log.Println("Conexão ao banco de dados estabelecida com sucesso")

    return DB, nil
}

