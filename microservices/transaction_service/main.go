package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "math/rand"
    "io/ioutil"
)

type TransactionStatus string

const (
    StatusPending   TransactionStatus = "pending"
    StatusSuccess   TransactionStatus = "success"
    StatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
    ID        uint              `json:"id" gorm:"primaryKey"`
    UserID    uint              `json:"user_id"`
    OfferID   uint              `json:"offer_id"`
    Amount    float64           `json:"amount"`
    Status    TransactionStatus `json:"status"`
    CreatedAt time.Time         `json:"created_at"`
    UpdatedAt time.Time         `json:"updated_at"`
}

// Variável global para a conexão com o banco de dados
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

func simulatePaymentGateway() bool {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(100) < 90 // 90% de chance de sucesso
}

// Função para criar uma nova transação com o status "pending"
func createTransaction(w http.ResponseWriter, r *http.Request) {
    userIDStr := r.URL.Query().Get("user_id")
    offerIDStr := r.URL.Query().Get("offer_id")
    amountStr := r.URL.Query().Get("amount")

    userID, _ := strconv.ParseUint(userIDStr, 10, 64)
    offerID, _ := strconv.ParseUint(offerIDStr, 10, 64)
    amount, _ := strconv.ParseFloat(amountStr, 64)

    // Criar nova transação
    transaction := Transaction{
        UserID:    uint(userID),
        OfferID:   uint(offerID),
        Amount:    amount,
        Status:    StatusPending,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := DB.Create(&transaction).Error; err != nil {
        log.Printf("Failed to create transaction: %v", err)
        http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
        return
    }

    log.Printf("Transaction created: %+v", transaction)
    json.NewEncoder(w).Encode(transaction)
}

func updateOfferUnits(offerID uint) error {
    // Montar a URL do microsserviço offer_update_service
    offerUpdateServiceURL := fmt.Sprintf("http://offer_update_service:8087/update-offer-units?id=%d", offerID)

    // Fazer a requisição HTTP GET ao offer_update_service
    resp, err := http.Get(offerUpdateServiceURL)
    if err != nil {
        return fmt.Errorf("failed to call offer_update_service: %v", err)
    }
    defer resp.Body.Close()

    // Verificar se a resposta foi bem-sucedida (status code 200)
    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := ioutil.ReadAll(resp.Body)
        bodyString := string(bodyBytes)
        return fmt.Errorf("offer_update_service returned non-OK status: %d, body: %s", resp.StatusCode, bodyString)
    }

    // Retornar nil se a chamada foi bem-sucedida
    return nil
}

// Função que processa o pagamento real utilizando transações de banco de dados
func processPayment(w http.ResponseWriter, r *http.Request) {
    transactionIDStr := r.URL.Query().Get("transaction_id")
    transactionID, _ := strconv.ParseUint(transactionIDStr, 10, 64)

    // Buscar a transação
    var transaction Transaction
    if err := DB.First(&transaction, transactionID).Error; err != nil {
        log.Printf("Transaction not found: %v", err)
        http.Error(w, "Transaction not found", http.StatusNotFound)
        return
    }

    // Se a transação já estiver marcada como success, evitar pagamento duplicado
    if transaction.Status == StatusSuccess {
        log.Printf("Transaction already successful: %v", transaction.ID)
        http.Error(w, "Transaction already successful", http.StatusBadRequest)
        return
    }

    // Iniciar uma transação de banco de dados
    tx := DB.Begin() // Inicia uma transação
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback() // Rollback se ocorrer algum erro inesperado
        }
    }()

    // Simular chamada ao gateway de pagamento real (substituir isso pela chamada real)
    paymentSuccess := simulatePaymentGateway()

    if paymentSuccess {
        // Atualizar o status da transação para "success" primeiro
        transaction.Status = StatusSuccess
        transaction.UpdatedAt = time.Now()

        if err := tx.Save(&transaction).Error; err != nil {
            log.Printf("Failed to update transaction: %v", err)
            tx.Rollback() // Se falhar, desfaz a transação
            http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
            return
        }

        // Tentar chamar o microsserviço offer_update_service para atualizar as unidades da oferta
        log.Printf("Calling offer_update_service to update units for offer ID: %d", transaction.OfferID)
        if err := updateOfferUnits(transaction.OfferID); err != nil {
            log.Printf("Failed to update offer units: %v", err)
            tx.Rollback() // Se falhar, desfaz a transação
            http.Error(w, "Failed to update offer units", http.StatusInternalServerError)
            return
        }

        // Tudo correu bem, confirma a transação
        if err := tx.Commit().Error; err != nil {
            log.Printf("Failed to commit transaction: %v", err)
            http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
            return
        }

        log.Printf("Transaction successful and offer updated: %+v", transaction)
        json.NewEncoder(w).Encode(transaction)

    } else {
        // Se o pagamento falhar, faz rollback da transação
        transaction.Status = StatusFailed
        transaction.UpdatedAt = time.Now()
        if err := tx.Save(&transaction).Error; err != nil {
            log.Printf("Failed to update failed transaction: %v", err)
            tx.Rollback() // Se falhar, desfaz a transação
            http.Error(w, "Failed to update failed transaction", http.StatusInternalServerError)
            return
        }

        // Tudo correu bem, então faz o commit
        if err := tx.Commit().Error; err != nil {
            log.Printf("Failed to commit failed transaction: %v", err)
            http.Error(w, "Failed to commit failed transaction", http.StatusInternalServerError)
            return
        }

        log.Printf("Transaction failed: %+v", transaction)
        json.NewEncoder(w).Encode(transaction)
    }
}

func main() {
    // Inicializar conexão com o banco de dados
    InitDB()

    // Endpoints para criar transação e processar pagamento
    http.HandleFunc("/create-transaction", createTransaction)
    http.HandleFunc("/process-payment", processPayment)

    // Define um handler HTTP para verificar o status do serviço
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Transaction Service is running!")
    })

    // Define a porta na qual o microsserviço vai escutar
    port := ":8085"
    log.Printf("Transaction Service is up and running on port %s", port)

    // Iniciar o servidor HTTP
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
