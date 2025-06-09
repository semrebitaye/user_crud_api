package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"user_crud_api/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		err = db.QueryRow("INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id", u.FirstName, u.LastName, u.Email, string(hash)).Scan(&u.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func GetUserById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u models.User
		err := db.QueryRow("SELECT * FROm users WHERE id = $1", id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)
	}
}

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
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

		json.NewEncoder(w).Encode(users)
	}
}
