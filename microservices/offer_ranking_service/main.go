// package main

// import (
//     "fmt"
//     "log"
//     "net/http"
//     "os"
//     "time"
//     "gorm.io/driver/postgres"
//     "gorm.io/gorm"
// )

// type Offer struct {
//     ID            uint
//     Price         float64
//     Discount      float64
//     Commission    float64
//     Tariff        float64
//     StartDate     time.Time
//     EndDate       time.Time
//     PromoUnits    uint64
//     UnitsSold     uint64
//     AvailableUnits uint64
//     Enable        bool
//     R1            float64
//     R2            float64
//     R3            float64
//     R4            float64
//     YNorm         float64
//     CreatedAt     time.Time
//     UpdatedAt     time.Time
// }

// var DB *gorm.DB

// // Inicializar a conexão com o banco de dados
// func InitDB() {
//     dsn := fmt.Sprintf(
//         "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
//         os.Getenv("DB_HOST"),
//         os.Getenv("DB_USER"),
//         os.Getenv("DB_PASSWORD"),
//         os.Getenv("DB_NAME"),
//         os.Getenv("DB_PORT"),
//     )

//     var err error
//     DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//     if err != nil {
//         log.Fatalf("Failed to connect to database: %v", err)
//     }

//     log.Println("Connected to the database successfully")
// }

// // Função para verificar os valores da oferta
// func verificarValores(offer *Offer) {
//     if offer.Price <= 0 || offer.PromoUnits == 0 || offer.StartDate.IsZero() || offer.EndDate.IsZero() {
//         offer.Enable = false
//     } else {
//         offer.Enable = true
//     }
// }

// // Função para calcular variáveis intermediárias
// func calcularVariaveisIntermediarias(offer *Offer) {
    
//     offer.R1 = offer.Discount / offer.Price
//     offer.R2 = offer.Commission / offer.Price
//     offer.R3 = float64(offer.UnitsSold) / float64(offer.PromoUnits)

//     now := time.Now()
//     durationTotal := offer.EndDate.Sub(offer.StartDate)
//     durationRemaining := offer.EndDate.Sub(now)
//     offer.R4 = transfCountDateTime(durationRemaining) / transfCountDateTime(durationTotal)
// }

// // Função para calcular YNorm com base nas variáveis intermediárias
// func calcularY(a, b, c, d, r1, r2, r3, r4 float64) float64 {
//     return a*r1 + b*r2 + c*r3 + d*r4
// }

// // Função para calcular YNorm com base nas variáveis intermediárias
// func calcularAvailableUnits(offer.AvailableUnits) int64 {
//     return offer.PromoUnits - offer.UnitsSold
// }

// // Função para transformar diferença entre datas em número decimal
// func transfCountDateTime(duration time.Duration) float64 {
//     return duration.Hours() / 24 // Convertendo para dias
// }


// // Função que roda em intervalos para atualizar as ofertas no banco de dados
// func UpdateOffers() error {
//     var offers []Offer
//     if err := DB.Find(&offers).Error; err != nil {
//         return err
//     }

//     for i := range offers {
//         offer := &offers[i]
//         verificarValores(offer)
//         if offer.Enable {
//             calcularVariaveisIntermediarias(offer)
//             calcularAvailableUnits(offer)
//             offer.AvailableUnits = offer.PromoUnits - offer.UnitsSold
//             offer.YNorm = calcularY(1.0, 2.0, -1.0, 0.5, offer.R1, offer.R2, offer.R3, offer.R4)
//         }
//     }

//     return DB.Save(&offers).Error
// }

// // Função que será chamada regularmente para ordenar e atualizar as ofertas
// func RunRankingService() {
//     ticker := time.NewTicker(1 * time.Minute) // Intervalo reduzido para teste
//     defer ticker.Stop()

//     for {
//         select {
//         case <-ticker.C:
//             log.Println("Running offer ranking update...")
//             if err := UpdateOffers(); err != nil {
//                 log.Printf("Error updating offers: %v", err)
//             } else {
//                 log.Println("Offers updated successfully")
//             }
//         }
//     }
// }

// func main() {
//     InitDB()

//     go RunRankingService()

//     http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, "Offer Ranking Service is running!")
//     })

//     port := ":8088"
//     log.Printf("Offer Ranking Service is up and running on port %s", port)

//     if err := http.ListenAndServe(port, nil); err != nil {
//         log.Fatalf("Failed to start server: %v", err)
//     }
// }


package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Offer struct {
    ID            uint
    Price         float64
    Discount      float64
    Commission    float64
    Tariff        float64
    StartDate     time.Time
    EndDate       time.Time
    PromoUnits    uint64
    UnitsSold     uint64
    AvailableUnits uint64
    Enable        bool
    R1            float64
    R2            float64
    R3            float64
    R4            float64
    YNorm         float64
    CreatedAt     time.Time
    UpdatedAt     time.Time
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

// Função para verificar os valores da oferta
func verificarValores(offer *Offer) {
    if offer.Price <= 0 || offer.PromoUnits == 0 || offer.StartDate.IsZero() || offer.EndDate.IsZero() {
        offer.Enable = false
    } else {
        offer.Enable = true
    }
}

// Função para calcular variáveis intermediárias
func calcularVariaveisIntermediarias(offer *Offer) {
    offer.R1 = offer.Discount / offer.Price
    offer.R2 = offer.Commission / offer.Price
    offer.R3 = float64(offer.UnitsSold) / float64(offer.PromoUnits)

    now := time.Now()
    durationTotal := offer.EndDate.Sub(offer.StartDate)
    durationRemaining := offer.EndDate.Sub(now)
    offer.R4 = transfCountDateTime(durationRemaining) / transfCountDateTime(durationTotal)
}

// Função para calcular YNorm com base nas variáveis intermediárias
func calcularY(a, b, c, d, r1, r2, r3, r4 float64) float64 {
    return a*r1 + b*r2 + c*r3 + d*r4
}

// Função para transformar diferença entre datas em número decimal
func transfCountDateTime(duration time.Duration) float64 {
    return duration.Hours() / 24 // Convertendo para dias
}

// Função que roda em intervalos para atualizar as ofertas no banco de dados
func UpdateOffers() error {
    var offers []Offer
    if err := DB.Find(&offers).Error; err != nil {
        return err
    }

    for i := range offers {
        offer := &offers[i]
        verificarValores(offer)
        if offer.Enable {
            calcularVariaveisIntermediarias(offer)
            offer.AvailableUnits = offer.PromoUnits - offer.UnitsSold // Atualiza AvailableUnits
            // offer.YNorm = calcularY(1.0, 2.0, -1.0, 0.5, offer.R1, offer.R2, offer.R3, offer.R4)
            offer.YNorm = calcularY(1.0, 0.01, 0.01, 0.01, offer.R1, offer.R2, offer.R3, offer.R4)
        }
    }

    return DB.Save(&offers).Error
}

// Função que será chamada regularmente para ordenar e atualizar as ofertas
func RunRankingService() {
    ticker := time.NewTicker(1 * time.Minute) // Intervalo reduzido para teste
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            log.Println("Running offer ranking update...")
            if err := UpdateOffers(); err != nil {
                log.Printf("Error updating offers: %v", err)
            } else {
                log.Println("Offer Ranking updated successfully")
            }
        }
    }
}

func main() {
    InitDB()

    go RunRankingService()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Offer Ranking Service is running!")
    })

    port := ":8088"
    log.Printf("Offer Ranking Service is up and running on port %s", port)

    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
