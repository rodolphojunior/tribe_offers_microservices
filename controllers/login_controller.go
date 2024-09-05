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
//         if err := config.DB.Preload("Company").Where("id = ?", user.ID).First(&partner).Error; err != nil {
//             http.Error(w, "Failed to load partner data", http.StatusInternalServerError)
//             return
//         }

//         response := map[string]interface{}{
//             "token":   token,
//             "user":    partner.User,
//             "company": partner.Company,
//         }
//         json.NewEncoder(w).Encode(response)
//     } else {
//         var consumer models.Consumer
//         if err := config.DB.Where("id = ?", user.ID).First(&consumer).Error; err != nil {
//             http.Error(w, "Failed to load consumer data", http.StatusInternalServerError)
//             return
//         }

//         response := map[string]interface{}{
//             "token": token,
//             "user":  consumer.User,
//         }
//         json.NewEncoder(w).Encode(response)
//     }
// }

package controllers

import (
    "encoding/json"
    "net/http"
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/models"
    "tribo_ofertas_backend/services"
    "tribo_ofertas_backend/utils"
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

    token, err := utils.GenerateJWT(user.ID, user.Role)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    if user.Role == "partner" {
        var partner models.Partner
        if err := config.DB.Preload("Company").Where("user_id = ?", user.ID).First(&partner).Error; err != nil {
            http.Error(w, "Failed to load partner data", http.StatusInternalServerError)
            return
        }

        response := map[string]interface{}{
            "token":   token,
            "user":    user,
            "company": partner.Company,
        }
        json.NewEncoder(w).Encode(response)
    } else if user.Role == "consumer" {
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
