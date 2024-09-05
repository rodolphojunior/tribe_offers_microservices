package controllers

import (
    "encoding/json"
    "net/http"
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/models"
    "tribo_ofertas_backend/services"
)

// func Register(w http.ResponseWriter, r *http.Request) {
//     var requestData struct {
//         User    models.User    `json:"user"`
//         Company *models.Company `json:"company,omitempty"`
//     }

//     if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
//         http.Error(w, "Invalid input", http.StatusBadRequest)
//         return
//     }

//     user := requestData.User

//     // Verificar se o email e senha foram realmente recebidos
//     if user.Email == "" || user.Password == "" {
//         http.Error(w, "email and password are required", http.StatusBadRequest)
//         return
//     }

//     // Se for um parceiro, associe com uma empresa
//     if user.Role == "partner" && requestData.Company != nil {
//         company := *requestData.Company

//         // Crie ou encontre a empresa baseada no CNPJ
//         if err := config.DB.FirstOrCreate(&company, models.Company{Cnpj: company.Cnpj}).Error; err != nil {
//             http.Error(w, "Failed to create or find company", http.StatusInternalServerError)
//             return
//         }

//         // Associar a empresa ao usuário antes de salvar o usuário
//         user.CompanyID = &company.ID
//     }

//     // Registre o usuário (agora com CompanyID preenchido)
//     if err := services.RegisterUser(&user, config.DB); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     // Crie o registro de Partner associando o User e a Company
//     if user.Role == "partner" && requestData.Company != nil {
//         partner := models.Partner{
//             UserID:    user.ID,
//             CompanyID: *user.CompanyID,
//         }

//         if err := config.DB.Create(&partner).Error; err != nil {
//             http.Error(w, "Failed to create partner", http.StatusInternalServerError)
//             return
//         }
//     }

//     // Se for um consumidor, crie o registro em `consumers`
//     if user.Role == "consumer" {
//         consumer := models.Consumer{
//             UserID: user.ID,
//         }
//         if err := config.DB.Create(&consumer).Error; err != nil {
//             http.Error(w, "Failed to create consumer", http.StatusInternalServerError)
//             return
//         }
//     }

//     w.WriteHeader(http.StatusCreated)
//     json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
// }

func Register(w http.ResponseWriter, r *http.Request) {
    var requestData struct {
        User    models.User    `json:"user"`
        Company *models.Company `json:"company,omitempty"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "Invalid input here", http.StatusBadRequest)
        return
    }

    user := requestData.User

    // Verificar se o email e senha foram realmente recebidos
    if user.Email == "" || user.Password == "" {
        http.Error(w, "email and password are required", http.StatusBadRequest)
        return
    }

    // Registre o usuário primeiro
    if err := services.RegisterUser(&user, config.DB); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Se for um parceiro, associe com uma empresa
    if user.Role == "partner" && requestData.Company != nil {
        company := *requestData.Company

        // Crie ou encontre a empresa baseada no CNPJ
        if err := config.DB.FirstOrCreate(&company, models.Company{Cnpj: company.Cnpj}).Error; err != nil {
            http.Error(w, "Failed to create or find company", http.StatusInternalServerError)
            return
        }

        // Crie o registro de Partner associando o User e a Company
        partner := models.Partner{
            UserID:    user.ID,
            CompanyID: company.ID,
        }

        if err := config.DB.Create(&partner).Error; err != nil {
            http.Error(w, "Failed to create partner", http.StatusInternalServerError)
            return
        }
    }

    // Se for um consumidor, crie o registro em `consumers`
    if user.Role == "consumer" {
        consumer := models.Consumer{
            UserID: user.ID,
        }
        if err := config.DB.Create(&consumer).Error; err != nil {
            http.Error(w, "Failed to create consumer", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}
