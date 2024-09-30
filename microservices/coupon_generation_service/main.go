package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "math/rand"  // Adiciona esta linha para importar o pacote correto
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type TransactionStatus string

const (
    StatusPending   TransactionStatus = "pending"
    StatusSuccess   TransactionStatus = "success"
    StatusFailed    TransactionStatus = "failed"
)

// Estruturas dos dados
type Offer struct {
    ID            uint
    Price         float64
    Title         string
    Description   string
    Discount      float64
    StartDate     time.Time
    EndDate       time.Time
    PromoUnits    uint
    UnitsSold     uint
    AvailableUnits uint
    Enable        bool
    CompanyID     uint
    Company       Company
}

type User struct {
    ID        uint
    FullName  string
    Cpf       string
}

type Transaction struct {
    ID            uint
    Amount        float64
    Status        TransactionStatus // success only
}

type Company struct {
    ID              uint
    Cnpj            string
    CompanyName     string
    Description     string
    TradeName       string
    Cep             string
    Address         string
    City            string
    State           string
    PhoneNumber     string
    Email           string
}

type Coupon struct {
    ID              uint
    Code            string
    Status          string
    Paid            bool
    OfferID         uint
    Offer           Offer
    UserID          uint
    User            User
    TransactionID   uint       
    Transaction     Transaction
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type CouponRequest struct {
    UserID        int `json:"user_id"`
    OfferID       int `json:"offer_id"`
    TransactionID int `json:"transaction_id"`
}

var DB *gorm.DB

// Inicializar a conexão com o banco de dados
func InitDB() {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    log.Println("Connected to the database successfully")
}

// Função para gerar o código do cupom
func generateCouponCode() string {
    timestamp := time.Now().Unix()
    return fmt.Sprintf("%d%s", timestamp, randomString(8))
}

// Função para gerar uma string aleatória
func randomString(n int) string {
    letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    s := make([]rune, n)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}

// Função para criar um cupom e retornar o JSON
func createCouponHandler(w http.ResponseWriter, r *http.Request) {

    // Verifica se o método da requisição é POST
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Decodifica o JSON da requisição para obter userID, offerID e transactionID
    var couponRequest CouponRequest
    if err := json.NewDecoder(r.Body).Decode(&couponRequest); err != nil {
        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
        return
    }

    userID := couponRequest.UserID
    offerID := couponRequest.OfferID
    transactionID := couponRequest.TransactionID

    // Buscar a oferta e verificar disponibilidade
    var offer Offer
    if err := DB.Preload("Company").First(&offer, offerID).Error; err != nil {
        http.Error(w, "Offer not found", http.StatusNotFound)
        return
    }
    if !offer.Enable || offer.AvailableUnits <= 0 {
        http.Error(w, "Offer not available", http.StatusBadRequest)
        return
    }

    // Buscar o usuário
    var user User
    if err := DB.First(&user, userID).Error; err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Buscar a transation
    var transaction Transaction
    if err := DB.First(&transaction, transactionID).Error; err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        return
    }
    if transaction.Status != StatusSuccess {
        http.Error(w, "Transacton not available", http.StatusBadRequest)
        return
    }

    // Criar o cupom
    coupon := Coupon{
        Code:           generateCouponCode(),
        Status:         "active",
        Paid:           true,
        OfferID:        offer.ID,
        Offer:          offer,
        UserID:         user.ID,
        User:           user,
        TransactionID:  transaction.ID,       
        Transaction:    transaction,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }

    // Salvar o cupom no banco de dados
    if err := DB.Create(&coupon).Error; err != nil {
        http.Error(w, "Failed to create coupon", http.StatusInternalServerError)
        return
    }

    // Retornar o JSON com os dados do cupom, oferta, usuário e empresa
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(coupon)
}

func main() {
    // Iniciar o banco de dados
    InitDB()

    // Definir o handler para gerar cupons
    http.HandleFunc("/generate-coupon", createCouponHandler)

    // Define um handler HTTP para verificar o status do serviço
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Coupon Generation Service is running!")
    })

    // Definir a porta para o microsserviço
    port := ":8084"
    log.Printf("Coupon Generation Service is up and running on port %s", port)

    // Iniciar o servidor HTTP
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

