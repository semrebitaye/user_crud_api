package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"user_crud_api/models"
)

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		json.NewDecoder(r.Body).Decode(&u)

		err := db.QueryRow("INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2) RETURNING id", u.FirstName, u.LastName, u.Email, u.Password).Scan(&u.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		json.NewEncoder(w).Encode(u)
	}
}
