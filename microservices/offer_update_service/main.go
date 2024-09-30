package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "time"
)

// Estrutura da oferta (Offer)
type Offer struct {
    ID              uint
    PromoUnits      uint64
    UnitsSold       uint64
    AvailableUnits  uint64
    Enable          bool 
    EndDate         time.Time   // Adicionar o campo EndDate
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

// Função para verificar se a oferta está disponível
func checkOfferAvailability(offer *Offer) bool {
    now := time.Now()

    // Verificar se a oferta está habilitada, tem unidades disponíveis e não expirou
    if offer.Enable && offer.AvailableUnits > 0 && offer.EndDate.After(now) {
        return true
    }

    // Se a oferta expirou, desabilitar a oferta e atualizá-la no banco de dados
    if offer.EndDate.Before(now) || offer.AvailableUnits == 0 {
        offer.Enable = false
        DB.Save(offer) // Atualizar a oferta no banco de dados
    }
    return false
}

// Função para atualizar os dados da oferta (vender unidades)
func updateOfferUnits(offer *Offer) error {
    // Verificar disponibilidade antes de atualizar
    if !checkOfferAvailability(offer) {
        return fmt.Errorf("offer is not available or has expired")
    }

    // Atualizar unidades
    offer.UnitsSold++
    offer.AvailableUnits--

    // Verificar disponibilidade após a venda
    if !checkOfferAvailability(offer) {
        return fmt.Errorf("offer has expired or is no longer available after update")
    }

    // Atualizar oferta no banco de dados
    return DB.Save(offer).Error
}

// Função que busca uma oferta pelo ID no banco de dados
func getOfferByID(id uint) (Offer, error) {
    var offer Offer
    if err := DB.First(&offer, id).Error; err != nil {
        return Offer{}, err
    }
    return offer, nil
}

// Handler para verificar a disponibilidade da oferta
func checkOfferHandler(w http.ResponseWriter, r *http.Request) {
    // Extrair o ID da oferta da URL
    offerIDStr := r.URL.Query().Get("id")
    offerID, err := strconv.ParseUint(offerIDStr, 10, 32)
    if err != nil {
        http.Error(w, "Invalid offer ID", http.StatusBadRequest)
        return
    }

    // Buscar a oferta no banco de dados
    offer, err := getOfferByID(uint(offerID))
    if err != nil {
        http.Error(w, "Offer not found", http.StatusNotFound)
        return
    }

    // Verificar a disponibilidade da oferta
    available := checkOfferAvailability(&offer)

    // Retornar o status da oferta
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":        available,
        "available":     offer.AvailableUnits,
        "units_sold":    offer.UnitsSold,
    })
}

// Handler para atualizar as unidades vendidas da oferta
func updateOfferUnitsHandler(w http.ResponseWriter, r *http.Request) {
    // Extrair o ID da oferta da URL
    offerIDStr := r.URL.Query().Get("id")
    offerID, err := strconv.ParseUint(offerIDStr, 10, 32)
    if err != nil {
        http.Error(w, "Invalid offer ID", http.StatusBadRequest)
        return
    }

    // Buscar a oferta no banco de dados
    offer, err := getOfferByID(uint(offerID))
    if err != nil {
        http.Error(w, "Offer not found", http.StatusNotFound)
        return
    }

    // Atualizar as unidades vendidas da oferta
    if err := updateOfferUnits(&offer); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Retornar o status da oferta atualizada
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":        "updated",
        "available":     offer.AvailableUnits,
        "units_sold":    offer.UnitsSold,
    })
}

func main() {
    // Inicializar a conexão com o banco de dados
    InitDB()

    // Definir o endpoint para verificar a oferta
    http.HandleFunc("/check-offer", checkOfferHandler)

    // Definir o endpoint para atualizar as unidades
    http.HandleFunc("/update-offer-units", updateOfferUnitsHandler)

    // Define um handler HTTP para verificar o status do serviço
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Offer Update Service is running!")
    })

    // Definir a porta para o microsserviço
    port := ":8087"
    log.Printf("Offer Update Service is up and running on port %s", port)

    // Iniciar o servidor HTTP
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
