package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"user_crud_api/models"

	"github.com/gorilla/mux"
)

func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		json.NewDecoder(r.Body).Decode(&u)

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE users SET name = $1, EMAIL = $2 WHERE id = $3", u.FirstName, u.LastName, u.Email, id, u.Password)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}
