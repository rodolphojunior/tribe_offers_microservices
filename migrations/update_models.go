package migrations

import (
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/models"
    "gorm.io/gorm"
)

func UpdateModels(db *gorm.DB) error {
    // Migrando as alterações nos modelos
    // err := db.AutoMigrate(&models.Company{}, &models.User{}, &models.Offer{})

	// err := db.AutoMigrate(&models.Offer{})
    // if err != nil {
    //     return err
    // }

    // Aqui você pode adicionar outras migrações específicas, como adição ou remoção de colunas
 
    // Exemplo: Remover uma coluna específica
    // if err := db.Migrator().DropColumn(&models.Offer{}, "avaiable_units"); err != nil {
    //     return err
    // }

    return nil
}

func RunMigrations() {
    db := config.DB
    if err := UpdateModels(db); err != nil {
        panic("Failed to update models: " + err.Error())
    }
}
