package controllers

import (
    "encoding/json"
    "net/http"
    "tribo_ofertas_backend/config"
    "tribo_ofertas_backend/models"
    "tribo_ofertas_backend/services"
    "tribo_ofertas_backend/utils"
    "errors"
)

// Login funcional
// Função Login que inclui consumers e partners
// func Login(w http.ResponseWriter, r *http.Request) {
//     var creds models.User
//     if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
//         http.Error(w, "Invalid input", http.StatusBadRequest)
//         return
//     }

//     // Autenticar o usuário usando o email e senha
//     user, err := services.AuthenticateUser(creds.Email, creds.Password, config.DB)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusUnauthorized)
//         return
//     }

//     // Verificar se o usuário está habilitado
//     if !user.Enabled {
//         http.Error(w, "User is disabled", http.StatusForbidden)
//         return
//     }

//     // Verificar o papel do usuário
//     if user.Role != "partner" && user.Role != "consumer" {
//         http.Error(w, "Unauthorized", http.StatusUnauthorized)
//         return
//     }

//     // Opcional: Retornar informações da Company associada, se necessário
//     config.DB.Preload("Company").Find(&user)
//     response := map[string]interface{}{
//         "token":   "jwt_token", // Gerar e retornar o JWT token aqui
//         "user":    user,
//         "company": user.Company,
//     }

//     json.NewEncoder(w).Encode(response)
// }

// Novo login experimental:
func Login(w http.ResponseWriter, r *http.Request) {
    var creds models.User
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Autenticar o usuário usando o email e senha
    user, err := services.AuthenticateUser(creds.Email, creds.Password, config.DB)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // Verificar se o usuário está habilitado
    if !user.Enabled {
        http.Error(w, "User is disabled", http.StatusForbidden)
        return
    }

    // Verificar o papel do usuário
    if user.Role != "partner" && user.Role != "consumer" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Gerar o JWT para o usuário autenticado
    token, err := utils.GenerateJWT(user.ID, user.Role)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Opcional: Retornar informações da Company associada, se necessário
    config.DB.Preload("Company").Find(&user)
    
    response := map[string]interface{}{
        "token":   token, // Retorna o JWT gerado
        "user":    user,
        "company": user.Company,
    }

    json.NewEncoder(w).Encode(response)
}


func validateUser(user models.User) error {
    if user.Email == "" || user.Password == "" {
        return errors.New("email and password are required")
    }

    // Verifique se o email ou username já existe no banco de dados
    var existingUser models.User
    if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
        return errors.New("user already exists with this email")
    }

    if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
        return errors.New("user already exists with this username")
    }

    return nil
}

// Função que Register inclui consumers e partners
func Register(w http.ResponseWriter, r *http.Request) {
    var requestData struct {
        User    models.User    `json:"user"`
        Company *models.Company `json:"company,omitempty"`
    }

    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    user := requestData.User
    company := requestData.Company

    if user.Email == "" || user.Password == "" {
        http.Error(w, "email and password are required", http.StatusBadRequest)
        return
    }

    if user.Role == "partner" && company != nil {
        var existingCompany models.Company
        if err := config.DB.Where("cnpj = ?", company.Cnpj).First(&existingCompany).Error; err == nil {
            user.CompanyID = &existingCompany.ID
        } else {
            if err := config.DB.Create(&company).Error; err != nil {
                http.Error(w, "Failed to create company", http.StatusInternalServerError)
                return
            }
            user.CompanyID = &company.ID
        }
    } else {
        user.CompanyID = nil // CompanyID deve ser nulo para consumers
    }

    if err := validateUser(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }
    user.Password = hashedPassword

    if err := config.DB.Create(&user).Error; err != nil {
        http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}
