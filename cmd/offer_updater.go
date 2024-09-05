package main

import (
    "log"
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/services"
    "time"
)

func main() {
    // Inicialize o banco de dados
    db, err := config.InitDB()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    for {
        // Atualize as ofertas
        err := services.UpdateOffers(db)
        if err != nil {
            log.Println("Error updating offers:", err)
        } else {
            log.Println("Offers updated successfully")
        }

        // Espera 24 horas até a próxima execução
        time.Sleep(24 * time.Hour)
    }
}


