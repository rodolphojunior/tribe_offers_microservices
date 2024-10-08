package routes

import (
    "github.com/gorilla/mux"
    "tribo_ofertas_backend/controllers"
    "tribo_ofertas_backend/middleware"
)

func InitRoutes() *mux.Router {
    router := mux.NewRouter()

    // Rota pública para retornar as ofertas válidas
    router.HandleFunc("/", controllers.GetPublicOffers).Methods("GET")

    // Rotas de autenticação
    router.HandleFunc("/register", controllers.Register).Methods("POST")
    router.HandleFunc("/login", controllers.Login).Methods("POST")

    // Rotas para consumidores (autenticados)
    consumersRouter := router.PathPrefix("/api/consumers").Subrouter()
    consumersRouter.Use(middleware.AuthMiddleware("consumer"))
    consumersRouter.HandleFunc("/offers", controllers.GetConsumerOffers).Methods("GET")

    // Rotas para parceiros
    partnersRouter := router.PathPrefix("/api/partners").Subrouter()
    partnersRouter.Use(middleware.AuthMiddleware("partner"))
    partnersRouter.HandleFunc("/offers", controllers.ManageOffers).Methods("POST")
    partnersRouter.HandleFunc("/my-offers", controllers.GetPartnerOffers).Methods("GET")

    return router
}



// package routes

// import (
//     "net/http"
//     "github.com/gorilla/mux"
//     "tribo_ofertas_backend/controllers"
//     "tribo_ofertas_backend/middleware"
// )

// func InitRoutes() *mux.Router {
//     router := mux.NewRouter()

//     // Rota de boas-vindas
//     router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//         w.Write([]byte("Bem-vindo à Tribo Ofertas API"))
//     }).Methods("GET")

//     // Rota de registro e login
//     router.HandleFunc("/register", controllers.Register).Methods("POST")
//     router.HandleFunc("/login", controllers.Login).Methods("POST")

//     // Rotas para consumidores
//     consumersRouter := router.PathPrefix("/api/consumers").Subrouter()
//     consumersRouter.Use(middleware.AuthMiddleware("consumer"))
//     consumersRouter.HandleFunc("/offers", controllers.GetConsumerOffers).Methods("GET")

//     // Rotas para parceiros
//     partnersRouter := router.PathPrefix("/api/partners").Subrouter()
//     partnersRouter.Use(middleware.AuthMiddleware("partner"))
//     partnersRouter.HandleFunc("/offers", controllers.ManageOffers).Methods("POST")

//     return router
// }


