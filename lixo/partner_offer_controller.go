package controllers

import (
    "net/http"
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/models"
    "encoding/json"
    "log" // Importa o pacote log
)

func GetPartnerOffers(w http.ResponseWriter, r *http.Request) {
    // Obter o userID do contexto
    userID, ok := r.Context().Value("user_id").(uint)
    if !ok {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var partnerAdmin models.PartnerAdmin
    if err := config.DB.Where("user_id = ?", user.ID).First(&partnerAdmin).Error; err != nil {
        http.Error(w, "Failed to load partner admin data", http.StatusInternalServerError)
        return
    }

    // Buscar as ofertas associadas à empresa (Company) do parceiro
    var offers []models.Offer
    if err := config.DB.Preload("Company").Where("company_id = ?", partner.CompanyID).Find(&offers).Error; err != nil {
        log.Printf("Error loading offers for company %d: %v", partner.CompanyID, err)
        http.Error(w, "Failed to load offers", http.StatusInternalServerError)
        return
    }

    // Registrar as ofertas carregadas nos logs
    log.Printf("Loaded offers for user partiner %d (partner %d, company %d): %+v", userID, partner.ID, partner.CompanyID, offers)

    // Enviar a resposta JSON com as ofertas
    json.NewEncoder(w).Encode(offers)
}

func ManageOffers(w http.ResponseWriter, r *http.Request) {
    var offer models.Offer

    // Decodificar o JSON recebido para o struct Offer
    if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Verificar se a empresa existe
    var company models.Company
    if err := config.DB.First(&company, offer.CompanyID).Error; err != nil {
        http.Error(w, "Company not found", http.StatusBadRequest)
        return
    }

    // Salvar a nova oferta no banco de dados
    if err := config.DB.Create(&offer).Error; err != nil {
        log.Println("Failed to create offer:", err)
        http.Error(w, "Failed to create offer: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Retornar uma resposta de sucesso
    log.Println("Offer created successfully")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Offer created successfully"})
}

