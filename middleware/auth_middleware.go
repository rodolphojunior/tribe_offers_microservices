// package middleware

// import (
//     "net/http"
//     "strings"
//     "tribo_ofertas_backend/utils"
//     "github.com/golang-jwt/jwt/v4"
// 	"github.com/gorilla/mux"
// )

// func AuthMiddleware(requiredRole string) mux.MiddlewareFunc {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             authHeader := r.Header.Get("Authorization")
//             if authHeader == "" {
//                 http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
//                 return
//             }

//             tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//             claims := &utils.Claims{}

//             token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//                 return utils.JwtKey, nil
//             })

//             if err != nil || !token.Valid {
//                 http.Error(w, "Invalid token", http.StatusUnauthorized)
//                 return
//             }

//             if claims.Role != requiredRole {
//                 http.Error(w, "Forbidden", http.StatusForbidden)
//                 return
//             }

//             r.Header.Set("userID", string(claims.UserID))
//             next.ServeHTTP(w, r)
//         })
//     }
// }

// func AuthMiddleware(role string) mux.MiddlewareFunc {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             tokenString := r.Header.Get("Authorization")

//             if tokenString == "" {
//                 http.Error(w, "Authorization header is required", http.StatusUnauthorized)
//                 return
//             }

//             tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

//             token, err := utils.ValidateJWT(tokenString)
//             if err != nil || token == nil {
//                 http.Error(w, "Invalid token", http.StatusUnauthorized)
//                 return
//             }

//             claims, ok := token.Claims.(jwt.MapClaims)
//             if !ok || !token.Valid {
//                 http.Error(w, "Invalid token claims", http.StatusUnauthorized)
//                 return
//             }

//             userID, ok := claims["user_id"].(float64)
//             if !ok {
//                 http.Error(w, "Invalid user ID in token claims", http.StatusUnauthorized)
//                 return
//             }

//             userRole, ok := claims["role"].(string)
//             if !ok || userRole != role {
//                 http.Error(w, "Unauthorized role", http.StatusUnauthorized)
//                 return
//             }

//             ctx := context.WithValue(r.Context(), "user_id", uint(userID))
//             next.ServeHTTP(w, r.WithContext(ctx))
//         })
//     }
// }

package middleware

import (
    "strings"
    "tribo_ofertas_backend/utils"
    "context"
    "net/http"
    "github.com/gorilla/mux" // Importa o pacote mux
    "github.com/golang-jwt/jwt/v4" // Importa o pacote jwt
)

func AuthMiddleware(role string) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenString := r.Header.Get("Authorization")

            if tokenString == "" {
                http.Error(w, "Authorization header is required", http.StatusUnauthorized)
                return
            }

            tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

            token, err := utils.ValidateJWT(tokenString)
            if err != nil || token == nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok || !token.Valid {
                http.Error(w, "Invalid token claims", http.StatusUnauthorized)
                return
            }

            userID, ok := claims["user_id"].(float64)
            if !ok {
                http.Error(w, "Invalid user ID in token claims", http.StatusUnauthorized)
                return
            }

            userRole, ok := claims["role"].(string)
            if !ok || userRole != role {
                http.Error(w, "Unauthorized role", http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), "user_id", uint(userID))
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
