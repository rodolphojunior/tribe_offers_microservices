package services

import (
    "errors"
    "tribe_offers_microservices/models"
    "tribe_offers_microservices/utils"
    "gorm.io/gorm"
)

// Valida se os campos obrigatórios do usuário estão preenchidos
func validateUser(user models.User) error {
    if user.Email == "" || user.Password == "" {
        return errors.New("email and password are required")
    }
    return nil
}

// Verifica se o usuário já existe no banco de dados
func checkUserExists(email string, db *gorm.DB) error {
    var existingUser models.User
    if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
        return errors.New("user already exists")
    }
    return nil
}

// Gera o hash da senha
func hashPassword(password string) (string, error) {
    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        return "", errors.New("failed to hash password")
    }
    return hashedPassword, nil
}

// Função principal que orquestra o registro do usuário
func RegisterUser(user *models.User, db *gorm.DB) error {
    // Verificação dos campos obrigatórios
    if err := validateUser(*user); err != nil {
        return err
    }

    // Verificação se o usuário já existe
    if err := checkUserExists(user.Email, db); err != nil {
        return err
    }

    // Hash da senha
    hashedPassword, err := hashPassword(user.Password)
    if err != nil {
        return err
    }
    user.Password = hashedPassword

    // Criação do usuário no banco de dados
    if err := db.Create(user).Error; err != nil {
        return errors.New("failed to create user")
    }

    return nil
}



// Autentica o usuário verificando o email e a senha
func AuthenticateUser(email string, password string, db *gorm.DB) (models.User, error) {
    var user models.User

    // Verifica se o usuário existe no banco de dados
    if err := db.Where("email = ?", email).First(&user).Error; err != nil {
        return user, errors.New("user not found")
    }

    // Verifica se a senha está correta
    if !utils.CheckPasswordHash(password, user.Password) {
        return user, errors.New("incorrect password")
    }

    return user, nil
}
