package main

import (
	"log"
	"net/http"
	"user_crud_api/controllers"
	initializer "user_crud_api/initializers"

	"github.com/gorilla/mux"
)

// var db *sql.DB

func main() {
	initializer.LoadEnvVariable()
	var db = initializer.Connect()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/create", controllers.CreateUser(db)).Methods("POST")
	r.HandleFunc("/login", controllers.Login(db)).Methods("POST")
	r.HandleFunc("/get", controllers.GetUsers(db)).Methods("GET")
	r.HandleFunc("/get/{id}", controllers.GetUserById(db)).Methods("GET")
	r.HandleFunc("/update/{id}", controllers.UpdateUser(db)).Methods("PATCH")
	r.HandleFunc("/delete/{id}", controllers.DeleteUser(db)).Methods("DELETE")
	// user := &models.User{ID: 1, FirstName: "man", LastName: "manega", Email: "man@man", Password: "manegaga"}
	// fmt.Println(user)
	log.Fatal(http.ListenAndServe(":8080", r))

}
