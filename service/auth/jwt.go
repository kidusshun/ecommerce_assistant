package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kidusshun/ecom_bot/config"
	"github.com/kidusshun/ecom_bot/utils"
)


var jwtSecret = []byte(config.JWTEnvs.JWTSecret)

func GenerateJWT(userEmail string) (string, error) {
	claims := jwt.MapClaims{
		"email": userEmail,         // Add user information to claims
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expiration (24 hours)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}



func CheckBearerToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("error 1")
            http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
            return
        }

        rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		tokenObj, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !tokenObj.Valid {
			log.Println("error 2",err)
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		claims, ok := tokenObj.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("error 3")
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		userEmail, ok := claims["email"].(string)
		if !ok {
			log.Println("error 4", err)
			utils.WriteError(w, http.StatusUnauthorized, errors.New("No email in token"))
			return
		}

		ctx := context.WithValue(r.Context(), "userEmail", userEmail)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
    })
}