package controllers

import (
    "net/http"
    // "time"
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/models"
    "encoding/json"
)

// // Função compartilhada para buscar e ordenar ofertas
// func fetchAndOrderOffers() ([]models.Offer, error) {
//     var offers []models.Offer
//     // Buscar e ordenar as ofertas por YNorm, assumindo que "enable" já está calculado
//     if err := config.DB.Where("enable = ?", true).Order("y_norm DESC").Find(&offers).Error; err != nil {
//         return nil, err
//     }
//     return offers, nil
// }

// Função compartilhada para buscar e ordenar ofertas
func fetchAndOrderOffers() ([]models.Offer, error) {
    var offers []models.Offer
    // Buscar e ordenar as ofertas por YNorm, assumindo que "enable" já está calculado
    if err := config.DB.Where("enable = ?", true).
        Order("y_norm DESC").
        Preload("Company").  // Preload para carregar os dados da empresa
        Preload("Coupons").  // Preload para carregar os dados de cupons
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


// func GetPublicOffers(w http.ResponseWriter, r *http.Request) {
// 	var offers []models.Offer
// 	if err := config.DB.Where("enable = ?", true).Order("y_norm DESC").Find(&offers).Error; err != nil {
// 		http.Error(w, "Failed to load offers", http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(offers)
// }


// func GetValidOffers(w http.ResponseWriter, r *http.Request) {
//     var offers []models.Offer

//     // Obtenha a data atual
//     currentDate := time.Now()

//     // Carregue ofertas válidas que estão dentro do prazo
//     if err := config.DB.Preload("Company").Where("start_date <= ? AND end_date >= ?", currentDate, currentDate).Find(&offers).Error; err != nil {
//         http.Error(w, "Failed to load offers", http.StatusInternalServerError)
//         return
//     }

//     // Verifique se encontrou alguma oferta
//     if len(offers) == 0 {
//         w.WriteHeader(http.StatusNoContent)
//         return
//     }

//     // Retorna as ofertas como JSON
//     json.NewEncoder(w).Encode(offers)
// }

