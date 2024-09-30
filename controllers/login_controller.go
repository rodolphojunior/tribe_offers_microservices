// package controllers

// import (
//     "encoding/json"
//     "net/http"
//     "tribo_ofertas_backend/config"
//     "tribo_ofertas_backend/models"
//     "tribo_ofertas_backend/services"
//     "tribo_ofertas_backend/utils"
// )

// func Login(w http.ResponseWriter, r *http.Request) {
//     var creds models.User
//     if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
//         http.Error(w, "Invalid input", http.StatusBadRequest)
//         return
//     }

//     user, err := services.AuthenticateUser(creds.Email, creds.Password, config.DB)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusUnauthorized)
//         return
//     }

//     if !user.Enabled {
//         http.Error(w, "User is disabled", http.StatusForbidden)
//         return
//     }

//     token, err := utils.GenerateJWT(user.ID, user.Role)
//     if err != nil {
//         http.Error(w, "Failed to generate token", http.StatusInternalServerError)
//         return
//     }

//     if user.Role == "partner" {
//         var partner models.Partner
//         if err := config.DB.Preload("Company").Where("user_id = ?", user.ID).First(&partner).Error; err != nil {
//             http.Error(w, "Failed to load partner data", http.StatusInternalServerError)
//             return
//         }

//         response := map[string]interface{}{
//             "token":   token,
//             "user":    user,
//             "company": partner.Company,
//         }
//         json.NewEncoder(w).Encode(response)
//     } else if user.Role == "consumer" {
//         var consumer models.Consumer
//         if err := config.DB.Where("user_id = ?", user.ID).First(&consumer).Error; err != nil {
//             http.Error(w, "Failed to load consumer data", http.StatusInternalServerError)
//             return
//         }

//         response := map[string]interface{}{
//             "token": token,
//             "user":  user,
//         }
//         json.NewEncoder(w).Encode(response)
//     } else {
//         http.Error(w, "Invalid role", http.StatusUnauthorized)
//     }
// }

// package controllers

// import (
//     "encoding/json"
//     "net/http"
//     "tribo_ofertas_backend/config"
//     "tribo_ofertas_backend/models"
//     "tribo_ofertas_backend/services"
//     "tribo_ofertas_backend/utils"
// )

// func Login(w http.ResponseWriter, r *http.Request) {
//     var creds models.User
//     if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
//         http.Error(w, "Invalid input", http.StatusBadRequest)
//         return
//     }

//     user, err := services.AuthenticateUser(creds.Email, creds.Password, config.DB)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusUnauthorized)
//         return
//     }

//     if !user.Enabled {
//         http.Error(w, "User is disabled", http.StatusForbidden)
//         return
//     }

//     token, err := utils.GenerateJWT(user.ID, string(user.Role))
//     if err != nil {
//         http.Error(w, "Failed to generate token", http.StatusInternalServerError)
//         return
//     }

//     // Lógica para PartnerAdmin
//     if user.Role == models.RoleAdmin {
//         var partnerAdmin models.PartnerAdmin
//         if err := config.DB.Preload("Company").Where("user_id = ?", user.ID).First(&partnerAdmin).Error; err != nil {
//             http.Error(w, "Failed to load partner admin data", http.StatusInternalServerError)
//             return
//         }

//         response := map[string]interface{}{
//             "token":   token,
//             "user":    user,
//             "company": partnerAdmin.Company,
//         }
//         json.NewEncoder(w).Encode(response)

//     // Lógica para PartnerAssistant
//     } else if user.Role == models.RoleAssistant {
//         var partnerAssistant models.PartnerAssistant
//         if err := config.DB.Preload("Company").Where("user_id = ?", user.ID).First(&partnerAssistant).Error; err != nil {
//             http.Error(w, "Failed to load partner assistant data", http.StatusInternalServerError)
//             return
//         }

//         response := map[string]interface{}{
//             "token":   token,
//             "user":    user,
//             "company": partnerAssistant.Company,
//         }
//         json.NewEncoder(w).Encode(response)

//     // Lógica para Consumer
//     } else if user.Role == models.RoleConsumer {
//         var consumer models.Consumer
//         if err := config.DB.Where("user_id = ?", user.ID).First(&consumer).Error; err != nil {
//             http.Error(w, "Failed to load consumer data", http.StatusInternalServerError)
//             return
//         }

//         response := map[string]interface{}{
//             "token": token,
//             "user":  user,
//         }
//         json.NewEncoder(w).Encode(response)
//     } else {
//         http.Error(w, "Invalid role", http.StatusUnauthorized)
//     }
// }

package controllers

import (
    "encoding/json"
    "net/http"
    "tribe_offers_microservices/config"
    "tribe_offers_microservices/models"
    "tribe_offers_microservices/services"
    "tribe_offers_microservices/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
    var creds models.User
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    user, err := services.AuthenticateUser(creds.Email, creds.Password, config.DB)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    if !user.Enabled {
        http.Error(w, "User is disabled", http.StatusForbidden)
        return
    }

    token, err := utils.GenerateJWT(user.ID, string(user.Role))
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Lógica para PartnerAdmin
    if user.Role == models.RoleAdmin {
        var partnerAdmin models.PartnerAdmin
        if err := config.DB.Preload("Company").Where("user_id = ?", user.ID).First(&partnerAdmin).Error; err != nil {
            http.Error(w, "Failed to load partner admin data", http.StatusInternalServerError)
            return
        }

        // Retornar dados essenciais da empresa, sem os dados sensíveis
        response := map[string]interface{}{
            "token":   token,
            "user":    user,
            "company": map[string]interface{}{
                "id":          partnerAdmin.Company.ID,
                "cnpj":        partnerAdmin.Company.Cnpj,
                "company_name": partnerAdmin.Company.CompanyName,
                "description": partnerAdmin.Company.Description,
                "trade_name":  partnerAdmin.Company.TradeName,
                "address":     partnerAdmin.Company.Address,
                "city":        partnerAdmin.Company.City,
                "state":       partnerAdmin.Company.State,
                "cep":         partnerAdmin.Company.Cep,
                "phone_number": partnerAdmin.Company.PhoneNumber,
                "email":       partnerAdmin.Company.Email,
                "website":     partnerAdmin.Company.Website,
            },
        }
        json.NewEncoder(w).Encode(response)

    // Lógica para PartnerAssistant
    } else if user.Role == models.RoleAssistant {
        var partnerAssistant models.PartnerAssistant
        if err := config.DB.Preload("Company").Where("user_id = ?", user.ID).First(&partnerAssistant).Error; err != nil {
            http.Error(w, "Failed to load partner assistant data", http.StatusInternalServerError)
            return
        }

        // Retornar apenas os dados essenciais da empresa e adicionar PartnerAdminID
        response := map[string]interface{}{
            "token":   token,
            "user":    user,
            "company": map[string]interface{}{
                "id":          partnerAssistant.Company.ID,
                "company_name": partnerAssistant.Company.CompanyName,
            },
            "partner_admin_id": partnerAssistant.PartnerAdminID, // Adiciona o ID do PartnerAdmin
        }
        json.NewEncoder(w).Encode(response)

    // Lógica para Consumer
    } else if user.Role == models.RoleConsumer {
        var consumer models.Consumer
        if err := config.DB.Where("user_id = ?", user.ID).First(&consumer).Error; err != nil {
            http.Error(w, "Failed to load consumer data", http.StatusInternalServerError)
            return
        }

        response := map[string]interface{}{
            "token": token,
            "user":  user,
        }
        json.NewEncoder(w).Encode(response)
    } else {
        http.Error(w, "Invalid role", http.StatusUnauthorized)
    }
}
