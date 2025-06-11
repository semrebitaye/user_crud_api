package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"user_crud_api/models"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusNotFound)
			log.Fatalf("Failed to decode the uder")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusNotFound)
			log.Fatalf("Failed to hash the password")
		}

		err = db.QueryRow(context.Background(), "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id", u.FirstName, u.LastName, u.Email, string(hash)).Scan(&u.ID)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusNotFound)
			log.Fatalf("Failed to get the requested id")
		}
		w.Header().Set("Conten-Type", "applicstion/json")
		json.NewEncoder(w).Encode(u)
	}
}

func GetUsers(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(context.Background(), "SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusNotFound)
			log.Fatalf("Failed to get the user")
		}
		defer rows.Close()

		users := []models.User{}
		for rows.Next() {
			var u models.User
			if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password); err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUserById(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u models.User
		err := db.QueryRow(context.Background(), "SELECT * FROm users WHERE id = $1", id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func UpdateUser(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		json.NewDecoder(r.Body).Decode(&u)

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec(context.Background(), "UPDATE users SET name = $1, EMAIL = $2 WHERE id = $3", u.FirstName, u.LastName, u.Email, id, u.Password)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func DeleteUser(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("User Deleted")
	}
}

func Login(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Println("failed to decode the req body")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var user models.User

		err = db.QueryRow(context.Background(), "SELECT email, password FROM \"users\" WHERE email=$1", body.Email).Scan(&user.Email, &user.Password)

		if err != nil {
			log.Println("the requested email does not exist")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

		if err != nil {
			log.Println("Failed to compare the password")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//generate a jwt tocken
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("secret")))
		if err != nil {
			log.Println("Failed to get the token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenString)
	}
}
