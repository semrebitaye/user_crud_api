package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func Authentication(next http.Handler) http.Handler {
	//get the Bearer of the req body
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("No token provided")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		//Decode/validateit
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v\n", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("secret")), nil
		})
		// Check for parsing error or invalid token
		if err != nil || !token.Valid {
			log.Printf("Token parse error: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Safe type assertion for claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// check exp date
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				log.Println("token expired")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "claims", token.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
