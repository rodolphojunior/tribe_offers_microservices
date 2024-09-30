package main

import (
    "log"
    "net/http"
    "tribe_offers_microservices/config"
    "tribe_offers_microservices/migrations"
    "tribe_offers_microservices/routes"
    "tribe_offers_microservices/middleware"
)

func main() {
    // Inicialize o banco de dados
    config.InitDB()

    // Execute as migrações
    migrations.RunMigrations()

    // Inicia as rotas com mux.Router
    router := routes.InitRoutes()

    // Aplicar o middleware de CORS
    wrappedRouter := middleware.EnableCORS(router)

    // Substituir http.Handle por ListenAndServe diretamente com o router
    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", wrappedRouter))
}


// package main

// import (
//     "log"
//     "fmt"
//     "net/http"
//     "tribe_offers_microservices/config"
//     "tribe_offers_microservices/routes"
//     "tribe_offers_microservices/migrations" // Descomentar para migrações


// )

// func enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")  // Permitir todas as origens
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	enableCors(&w) // Habilitar CORS em cada requisição

// 	if r.Method == "OPTIONS" {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	// Lógica para lidar com o POST
// }

// func main() {
//     // Inicialize o banco de dados
//     config.InitDB()

//     // Execute as migrações
//     migrations.RunMigrations()

//     // Inicia as rotas
//     router := routes.InitRoutes()

// 	http.HandleFunc("/register", handler)
// 	http.ListenAndServe(":8080", nil)

//     log.Println("Server is running on port 8080")
//     log.Fatal(http.ListenAndServe(":8080", router))

//     http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, "Pong")
//     })
// }
 