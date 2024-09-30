// package controllers

// import (
//     "encoding/json"
//     "net/http"
//     "tribo_ofertas_backend/config"
//     "tribo_ofertas_backend/models"
//     "tribo_ofertas_backend/services"
// )

// func Register(w http.ResponseWriter, r *http.Request) {
//     var requestData struct {
//         User    models.User    `json:"user"`
//         Company *models.Company `json:"company,omitempty"`
//     }
    
//     if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
//         http.Error(w, "Invalid input here", http.StatusBadRequest)
//         return
//     }

//     user := requestData.User

//     // Verificar se o email e senha foram realmente recebidos
//     if user.Email == "" || user.Password == "" {
//         http.Error(w, "email and password are required", http.StatusBadRequest)
//         return
//     }

//     // Registre o usuário primeiro
//     if err := services.RegisterUser(&user, config.DB); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
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

//         // Crie o registro de Partner associando o User e a Company
//         partner := models.Partner{
//             UserID:    user.ID,
//             CompanyID: company.ID,
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

package controllers

import (
    "encoding/json"
    "log"
    "net/http"
    "tribe_offers_microservices/config"
    "tribe_offers_microservices/models"
    "tribe_offers_microservices/services"
    "io"
)

func Register(w http.ResponseWriter, r *http.Request) {
    var requestData struct {
        User            models.User     `json:"user"`
        Company         *models.Company `json:"company,omitempty"`
        PartnerAdminID  *uint           `json:"partner_admin_id,omitempty"`
    }

    log.Println("Recebendo dados de registro...")

    // Ler o corpo da requisição e registrar o conteúdo para debug
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Erro ao ler o corpo da requisição: %v", err)
        http.Error(w, "Erro ao ler o corpo da requisição", http.StatusInternalServerError)
        return
    }
    log.Println("Corpo da requisição recebido (string): ", string(bodyBytes))

    // Decodificar o JSON
    if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
        log.Printf("Erro ao decodificar input: %v", err)
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    log.Printf("Dados decodificados com sucesso: %+v", requestData)
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

    // Verificar se o papel é Admin e associar à empresa
    if user.Role == "admin" && requestData.Company != nil {
        company := *requestData.Company

        // Crie ou encontre a empresa baseada no CNPJ
        if err := config.DB.FirstOrCreate(&company, models.Company{Cnpj: company.Cnpj}).Error; err != nil {
            http.Error(w, "Failed to create or find company", http.StatusInternalServerError)
            return
        }

        // Crie o registro de PartnerAdmin associando o User e a Company
        partnerAdmin := models.PartnerAdmin{
            UserID:    user.ID,
            CompanyID: company.ID,
        }

        if err := config.DB.Create(&partnerAdmin).Error; err != nil {
            http.Error(w, "Failed to create partner admin", http.StatusInternalServerError)
            return
        }
    }

    // Se o papel for Assistant, associe o PartnerAssistant ao PartnerAdmin e à empresa
    if user.Role == "assistant" && requestData.Company != nil && requestData.PartnerAdminID != nil {
        company := *requestData.Company

        // Encontre o PartnerAdmin associado ao PartnerAssistant usando o PartnerAdminID
        var partnerAdmin models.PartnerAdmin
        if err := config.DB.Where("id = ?", *requestData.PartnerAdminID).First(&partnerAdmin).Error; err != nil {
            http.Error(w, "Failed to find partner admin", http.StatusInternalServerError)
            return
        }

        // Crie o registro de PartnerAssistant associando o User e a Company
        partnerAssistant := models.PartnerAssistant{
            UserID:       user.ID,
            PartnerAdminID: *requestData.PartnerAdminID,
            CompanyID:    company.ID,
        }

        if err := config.DB.Create(&partnerAssistant).Error; err != nil {
            http.Error(w, "Failed to create partner assistant", http.StatusInternalServerError)
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
