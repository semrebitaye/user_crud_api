package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"user_crud_api/models"

	"github.com/gorilla/mux"
)

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
