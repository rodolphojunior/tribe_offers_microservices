package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "github.com/skip2/go-qrcode"
)

// Estrutura para receber o JSON
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

// Função para criptografar dados usando AES
func encrypt(data []byte, passphrase string) ([]byte, error) {
    block, err := aes.NewCipher([]byte(passphrase))
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    _, err = io.ReadFull(rand.Reader, nonce)
    if err != nil {
        return nil, err
    }
    return gcm.Seal(nonce, nonce, data, nil), nil
}

// Função para gerar o QR Code criptografado e codificado em Base64
func generateQRCodeHandler(w http.ResponseWriter, r *http.Request) {
    // Decodificar o JSON recebido no corpo da requisição
    var coupon Coupon
    if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
        return
    }

    // Converter o cupom para JSON
    couponData, err := json.Marshal(coupon)
    if err != nil {
        http.Error(w, "Failed to marshal coupon data", http.StatusInternalServerError)
        return
    }

    // Criptografar os dados do cupom
    encryptedData, err := encrypt(couponData, "1234567890123456")
    if err != nil {
        http.Error(w, "Failed to encrypt data", http.StatusInternalServerError)
        return
    }

    // Codificar os dados criptografados em Base64
    encodedData := base64.StdEncoding.EncodeToString(encryptedData)

    // Ensure the output directory exists
    outputDir := "/app/output"
    if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
        log.Printf("Failed to create output directory: %v", err)
        http.Error(w, "Failed to create output directory", http.StatusInternalServerError)
        return
    }

    // Generate the QR Code
    qrCodePath := filepath.Join(outputDir, fmt.Sprintf("encrypted_coupon_%s.png", coupon.Code))
    log.Printf("Generating QR Code for coupon with Code: %s", coupon.Code)
    log.Printf("QR Code will be saved as: %s", qrCodePath)

    err = qrcode.WriteFile(encodedData, qrcode.Medium, 256, qrCodePath)
    if err != nil {
        log.Printf("Failed to generate QR code: %v", err)
        http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
        return
    }

    // Check if the file was created
    if _, err := os.Stat(qrCodePath); os.IsNotExist(err) {
        log.Printf("QR Code file was not created: %v", err)
        http.Error(w, "Failed to create QR code file", http.StatusInternalServerError)
        return
    }

    log.Printf("QR Code successfully generated and saved as %s", qrCodePath)
    w.Write([]byte(fmt.Sprintf("QR Code successfully generated and saved as %s", qrCodePath)))
}

func main() {
    http.HandleFunc("/generate-qrcode", generateQRCodeHandler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "QR Code Service is running!")
    })

    port := ":8086"
    log.Printf("QR Code Service is up and running on port %s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
