package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"user_crud_api/models"

	"github.com/gorilla/mux"
)

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

func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode("User Deleted")
	}
}

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
