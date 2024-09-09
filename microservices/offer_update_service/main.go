package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    // Define um handler que responde a requisição HTTP
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Offer Update Service is running!")
    })

    // Define a porta na qual o microsserviço vai escutar
    port := ":8087"

    log.Printf("Offer Update Service is up and running on port %s", port)

    // Inicia o servidor HTTP
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
