package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "time"
)

type Offer struct {
	ID            	uint      `json:"id" gorm:"primaryKey"`
	Price         	float64   `json:"price"`
    Title  		  	string    `json:"title"`
    Description   	string    `json:"description"`
	AverageRating 	float64   `json:"average_rating" gorm:"default:0.0"` // Valor padrão de 0.0
    CompanyID     	uint      `json:"company_id"`
}

type OfferEvaluation struct {
    ID          uint    `gorm:"primaryKey"`
    OfferID     uint    `json:"offer_id"`  // Referência para a oferta
    Offer       Offer   `gorm:"foreignKey:OfferID"`
    UserID      uint    `json:"user_id"`   // Referência para o cliente que avaliou
    CouponID    uint    `json:"coupon_id"` // Referência para o cupom resgatado
    Rating      float64 `json:"rating" gorm:"default:0.0"`    // Nota dada pelo cliente (exemplo: 1 a 5 estrelas)
    Comment     string  `json:"comment"`   // Comentário opcional do cliente
    CreatedAt   time.Time `json:"created_at"`
}

var DB *gorm.DB

// Função para inicializar o banco de dados
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
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    log.Println("Connected to the database successfully")
}

// Função para criar uma nova avaliação e atualizar a média da oferta
type OfferEvaluationStats struct {
    TotalRatings float64 `gorm:"column:total_ratings"`
    Count        int64   `gorm:"column:count"`
}

func CreateOfferEvaluation(db *gorm.DB, evaluation *OfferEvaluation) error {
    // Crie a nova avaliação no banco de dados
    if err := db.Create(evaluation).Error; err != nil {
        return err
    }

    // Verifique se existem outras avaliações para a oferta
    var count int64
    if err := db.Model(OfferEvaluation{}).Where("offer_id = ?", evaluation.OfferID).Count(&count).Error; err != nil {
        return err
    }

    // Apenas atualize a média de ratings se houver outras avaliações
    if count > 1 {
        var stats OfferEvaluationStats
        if err := db.Model(OfferEvaluation{}).
            Where("offer_id = ?", evaluation.OfferID).
            Select("AVG(rating) as total_ratings, COUNT(*) as count").
            Scan(&stats).Error; err != nil {
            return err
        }

        // Atualize a oferta com a nova média de rating
        return db.Model(&Offer{}).
            Where("id = ?", evaluation.OfferID).
            Update("average_rating", stats.TotalRatings).Error
    }

    return nil
}

// Handler para receber avaliações de ofertas
func CreateOfferEvaluationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var evaluation OfferEvaluation
    if err := json.NewDecoder(r.Body).Decode(&evaluation); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    log.Println("PASSOU POR AQUI")
    if err := CreateOfferEvaluation(DB, &evaluation); err != nil {
        http.Error(w, "Failed to create evaluation", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Evaluation created successfully"})
}

func main() {
    // Inicializar o banco de dados
    InitDB()

    // Define um handler que responde a requisição HTTP
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Coupon Evaluation Service is running!")
    })

    // Rota para criar avaliações
    http.HandleFunc("/create-evaluation", CreateOfferEvaluationHandler)

    // Define a porta na qual o microsserviço vai escutar
    port := ":8090"
    log.Printf("Coupon Evaluation Service is up and running on port %s", port)

    // Inicia o servidor HTTP
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
