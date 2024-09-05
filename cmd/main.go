package main

import (
    "log"
    "net/http"
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/routes"
    // "tribo_ofertas_backend/migrations" // Descomentar para migrações


)

func main() {
    // Inicialize o banco de dados
    config.InitDB()

    // Execute as migrações para os novos models
    // config.DB.AutoMigrate(&models.Company{}, &models.User{}, &models.Offer{})

    // Execute as migrações
    // migrations.RunMigrations()

    // Inicia as rotas
    router := routes.InitRoutes()

    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
