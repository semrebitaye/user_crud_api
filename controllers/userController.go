package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"user_crud_api/models"

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

		json.NewEncoder(w).Encode(u)
	}
}
