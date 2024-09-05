package controllers

// import (
//     "encoding/json"
//     "log"
//     "net/http"
//     "tribo_ofertas_backend/config"
//     "tribo_ofertas_backend/models"
//     // "github.com/golang-jwt/jwt/v4"
//     "github.com/dgrijalva/jwt-go"
//     // "github.com/gorilla/mux"
//     // "strconv"
//     "time"
// )


import (
    "net/http"
    // "time"
    // "github.com/gorilla/mux"
    // "github.com/golang-jwt/jwt/v4"
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/models"
    "encoding/json"
    "log" // Importa o pacote log
)

func GetPartnerOffers(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(uint)
    if !ok {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var partner models.Partner
    if err := config.DB.Where("user_id = ?", userID).First(&partner).Error; err != nil {
        http.Error(w, "Failed to load partner data", http.StatusInternalServerError)
        return
    }

    var offers []models.Offer
    if err := config.DB.Where("company_id = ?", partner.CompanyID).Find(&offers).Error; err != nil {
        http.Error(w, "Failed to load offers", http.StatusInternalServerError)
        return
    }

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

// func GetPartnerOffers(w http.ResponseWriter, r *http.Request) {
//     // Obtenha o token JWT do contexto
//     user := r.Context().Value("user").(*jwt.Token)
//     claims := user.Claims.(jwt.MapClaims)

//     // Extrair o user_id dos claims
//     userID := uint(claims["user_id"].(float64))

//     var partner models.Partner
//     if err := config.DB.Where("user_id = ?", userID).First(&partner).Error; err != nil {
//         http.Error(w, "Partner not found", http.StatusNotFound)
//         return
//     }

//     var offers []models.Offer
//     currentDate := time.Now()
//     if err := config.DB.Where("company_id = ? AND start_date <= ? AND end_date >= ?", partner.CompanyID, currentDate, currentDate).Find(&offers).Error; err != nil {
//         http.Error(w, "Failed to load offers", http.StatusInternalServerError)
//         return
//     }

//     json.NewEncoder(w).Encode(offers)
// }

// func GetPartnerOffers(w http.ResponseWriter, r *http.Request) {
//     // Obtenha o token JWT do contexto
//     user := r.Context().Value("user").(*jwt.Token)
//     claims := user.Claims.(jwt.MapClaims)

//     // Extrair o user_id dos claims
//     userID := uint(claims["user_id"].(float64))
//     log.Printf("Extracted user ID from JWT: %d", userID)

//     var partner models.Partner
//     if err := config.DB.Where("user_id = ?", userID).First(&partner).Error; err != nil {
//         log.Printf("Error finding partner: %v", err)
//         http.Error(w, "Partner not found", http.StatusNotFound)
//         return
//     }

//     var offers []models.Offer
//     currentDate := time.Now()
//     if err := config.DB.Where("company_id = ? AND start_date <= ? AND end_date >= ?", partner.CompanyID, currentDate, currentDate).Find(&offers).Error; err != nil {
//         log.Printf("Error loading offers: %v", err)
//         http.Error(w, "Failed to load offers", http.StatusInternalServerError)
//         return
//     }

//     log.Printf("Returning %d offers", len(offers))
//     json.NewEncoder(w).Encode(offers)
// }

// func GetPartnerOffers(w http.ResponseWriter, r *http.Request) {
//     userID, ok := r.Context().Value("user_id").(uint)
//     if !ok {
//         http.Error(w, "Invalid user ID", http.StatusBadRequest)
//         return
//     }

//     var partner models.Partner
//     if err := config.DB.Where("user_id = ?", userID).First(&partner).Error; err != nil {
//         http.Error(w, "Failed to load partner data", http.StatusInternalServerError)
//         return
//     }

//     var offers []models.Offer
//     if err := config.DB.Where("company_id = ?", partner.CompanyID).Find(&offers).Error; err != nil {
//         http.Error(w, "Failed to load offers", http.StatusInternalServerError)
//         return
//     }

//     json.NewEncoder(w).Encode(offers)
// }
