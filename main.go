package main

import (
	"database/sql"
	"log"
	"net/http"
	"user_crud_api/controllers"
	initializer "user_crud_api/initializers"

	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	initializer.LoadEnvVariable()
	initializer.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/update/{id}", controllers.UpdateUser(db)).Methods("PUT")
	r.HandleFunc("/delete/{id}", controllers.DeleteUser(db)).Methods("DELETE")

	// user := &models.User{ID: 1, FirstName: "man", LastName: "manega", Email: "man@man", Password: "manegaga"}
	// fmt.Println(user)
	log.Fatal(http.ListenAndServe(":8080", r))

}
