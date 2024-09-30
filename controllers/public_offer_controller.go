package controllers

import (
    "net/http"
    // "time"
    "tribe_offers_microservices/config"
    "tribe_offers_microservices/models"
    "encoding/json"
)

// Função compartilhada para buscar e ordenar ofertas
func fetchAndOrderOffers() ([]models.Offer, error) {
    var offers []models.Offer
    // Buscar e ordenar as ofertas por YNorm, assumindo que "enable" já está calculado
    if err := config.DB.Where("enable = ?", true).
        Order("y_norm DESC").
        Preload("Company").  // Preload para carregar os dados da empresa
        //Preload("Coupons").  // Preload para carregar os dados de cupons
        Preload("Photos").  // Preload para carregar os dados de fotos
        Find(&offers).Error; err != nil {
        return nil, err
    }
    return offers, nil
}

// Controller para ofertas públicas
func GetPublicOffers(w http.ResponseWriter, r *http.Request) {
    offers, err := fetchAndOrderOffers()
    if err != nil {
        http.Error(w, "Failed to load offers", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(offers)
}

// Controller para ofertas para consumidores autenticados
func GetConsumerOffers(w http.ResponseWriter, r *http.Request) {
    offers, err := fetchAndOrderOffers()
    if err != nil {
        http.Error(w, "Failed to load offers", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(offers)
}
