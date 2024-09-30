package main

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/makiuchi-d/gozxing"
    "github.com/makiuchi-d/gozxing/qrcode"
    "log"
    "net/http"
    "os"
    "image"
    _ "image/png"
    "github.com/nfnt/resize"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// Estrutura completa do Cupom
type Coupon struct {
    ID            uint   `json:"id"`
    Code          string `json:"code"`
    Status        string `json:"status"`
    Paid          bool   `json:"paid"`
    OfferID       uint   `json:"offer_id"`
    UserID        uint   `json:"user_id"`
    TransactionID uint   `json:"transaction_id"`
    CreatedAt     string `json:"created_at"`
    UpdatedAt     string `json:"updated_at"`
}

// Estrutura da Oferta
type Offer struct {
    ID          uint   `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Price       float64 `json:"price"`
    Discount    float64 `json:"discount"`
    StartDate   string `json:"start_date"`
    EndDate     string `json:"end_date"`
    CompanyID   uint   `json:"company_id"`
}

// Estrutura da Empresa
type Company struct {
    ID          uint   `json:"id"`
    CNPJ        string `json:"cnpj"`
    CompanyName string `json:"company_name"`
    Address     string `json:"address"`
    City        string `json:"city"`
    State       string `json:"state"`
}

// Estrutura do User
type User struct {
    ID       uint   `json:"id"`
    FullName string `json:"full_name"`
    Cpf      string `json:"cpf"`
    Email    string `json:"email"`
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

// Função para descriptografar dados usando AES
func decrypt(data []byte, passphrase string) ([]byte, error) {
    block, err := aes.NewCipher([]byte(passphrase))
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, fmt.Errorf("invalid data")
    }
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

// Função para buscar o cupom completo pelo ID
func getCouponByID(couponID uint) (Coupon, error) {
    var coupon Coupon
    if err := DB.First(&coupon, couponID).Error; err != nil {
        return coupon, err
    }
    return coupon, nil
}

// Função para buscar o usuário pelo ID
func getUserByID(userID uint) (User, error) {
    var user User
    if err := DB.First(&user, userID).Error; err != nil {
        return user, err
    }
    return user, nil
}

// Função para buscar a oferta pelo ID
func getOfferByID(offerID uint) (Offer, error) {
    var offer Offer
    if err := DB.First(&offer, offerID).Error; err != nil {
        return offer, err
    }
    return offer, nil
}

// Função para buscar a empresa pelo ID
func getCompanyByID(companyID uint) (Company, error) {
    var company Company
    if err := DB.First(&company, companyID).Error; err != nil {
        return company, err
    }
    return company, nil
}

// Handler para decodificar e validar o QR Code
func redeemCouponHandler(w http.ResponseWriter, r *http.Request) {
    fileName := r.URL.Query().Get("file")
    if fileName == "" {
        http.Error(w, "No file provided", http.StatusBadRequest)
        return
    }

    filePath := fmt.Sprintf("/app/output/%s", fileName)
    log.Printf("Trying to open file: %s", filePath)

    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("Failed to open file: %v", err)
        http.Error(w, "Failed to open file", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        log.Printf("Failed to decode image: %v", err)
        http.Error(w, "Failed to decode image", http.StatusInternalServerError)
        return
    }

    img = resize.Resize(256, 256, img, resize.Lanczos3)

    bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
    if err != nil {
        log.Printf("Failed to convert image to binary bitmap: %v", err)
        http.Error(w, "Failed to process image", http.StatusInternalServerError)
        return
    }

    qrReader := qrcode.NewQRCodeReader()
    result, err := qrReader.Decode(bitmap, nil)
    if err != nil {
        log.Printf("Failed to decode QR code: %v", err)
        http.Error(w, "Failed to decode QR code", http.StatusInternalServerError)
        return
    }

    encryptedData, err := base64.StdEncoding.DecodeString(result.GetText())
    if err != nil {
        log.Printf("Failed to decode QR code content: %v", err)
        http.Error(w, "Failed to decode QR code content", http.StatusBadRequest)
        return
    }

    decryptedData, err := decrypt(encryptedData, "1234567890123456")
    if err != nil {
        log.Printf("Failed to decrypt QR code: %v", err)
        http.Error(w, "Failed to decrypt QR code", http.StatusBadRequest)
        return
    }

    var coupon Coupon
    if err := json.Unmarshal(decryptedData, &coupon); err != nil {
        log.Printf("Failed to parse decrypted coupon data: %v", err)
        http.Error(w, "Failed to parse decrypted coupon data", http.StatusInternalServerError)
        return
    }

    log.Printf("Coupon redeemed: %+v", coupon)

    if coupon.Status != "active" {
        http.Error(w, "Coupon is not active", http.StatusBadRequest)
        return
    }

    // Buscar detalhes do cupom
    couponDetails, err := getCouponByID(coupon.ID)
    if err != nil {
        log.Printf("Failed to get coupon details: %v", err)
        http.Error(w, "Failed to get coupon details", http.StatusInternalServerError)
        return
    }

    // Buscar detalhes da oferta associada
    offer, err := getOfferByID(couponDetails.OfferID)
    if err != nil {
        log.Printf("Failed to get offer: %v", err)
        http.Error(w, "Failed to get offer", http.StatusInternalServerError)
        return
    }

    // Buscar detalhes da empresa associada
    company, err := getCompanyByID(offer.CompanyID)
    if err != nil {
        log.Printf("Failed to get company: %v", err)
        http.Error(w, "Failed to get company", http.StatusInternalServerError)
        return
    }

    // Buscar detalhes do usuário associado ao cupom
    user, err := getUserByID(couponDetails.UserID)
    if err != nil {
        log.Printf("Failed to get user: %v", err)
        http.Error(w, "Failed to get user", http.StatusInternalServerError)
        return
    }

    // Retornar os dados como resposta JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Coupon successfully redeemed!",
        "coupon":  couponDetails,
        "offer": map[string]interface{}{
            "id":          offer.ID,
            "title":       offer.Title,
            "description": offer.Description,
            "price":       offer.Price,
            "discount":    offer.Discount,
            "start_date":  offer.StartDate,
            "end_date":    offer.EndDate,
        },
        "company": map[string]interface{}{
            "cnpj":         company.CNPJ,
            "company_name": company.CompanyName,
            "address":      company.Address,
            "city":         company.City,
            "state":        company.State,
        },
        "user": map[string]interface{}{
            "full_name": user.FullName,
            "cpf":       user.Cpf,
            "email":     user.Email,
        },
    })
}

func main() {
    // Inicializar o banco de dados
    InitDB()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Coupon Redemption Service is running!")
    })

    http.HandleFunc("/redeem-coupon", redeemCouponHandler)

    port := ":8089"
    log.Printf("Coupon Redemption Service is up and running on port %s", port)

    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
