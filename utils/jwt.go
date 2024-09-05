package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
    "errors"
	
)

var JwtKey = []byte("your_secret_key")

type Claims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.StandardClaims
}

func GenerateJWT(userID uint, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}


var jwtKey = []byte("your_secret_key")

func ValidateJWT(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if token.Valid {
        return token, nil
    } else {
        return nil, errors.New("invalid token")
    }
}
