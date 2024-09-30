package migrations

import (
    "tribe_offers_microservices/config"
    "tribe_offers_microservices/models"
    "gorm.io/gorm"
    "log"
)

// Função que executa as migrações no banco de dados
func UpdateModels(db *gorm.DB) error {
    // AutoMigrate dos modelos
    err := db.AutoMigrate(
        &models.Offer{},
        &models.OfferEvaluation{},
    )
    if err != nil {
        return err
    }
    return nil
}

// Função que roda as migrações
func RunMigrations() {
    // Verifica se o banco de dados foi inicializado
    db := config.DB
    if db == nil {
        log.Fatal("Falha ao inicializar o banco de dados: instância DB é nula")
    }

    // Executa a função de migração
    if err := UpdateModels(db); err != nil {
        log.Fatalf("Falha ao atualizar os modelos: %v", err)
    }

    log.Println("Migrações concluídas com sucesso")
}

