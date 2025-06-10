package main

import (
	"context"
	"log"
	"net/http"
	"user_crud_api/controllers"
	initializer "user_crud_api/initializers"
	"user_crud_api/middleware"

	"github.com/gorilla/mux"
)

// var db *sql.DB

func main() {
	initializer.LoadEnvVariable()
	db := initializer.Connect()
	defer db.Close(context.Background())

	r := mux.NewRouter()
	r.HandleFunc("/users", controllers.CreateUser(db)).Methods("POST")
	r.HandleFunc("/login", controllers.Login(db)).Methods("POST")

	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.Authentication)
	api.HandleFunc("/users", controllers.GetUsers(db)).Methods("GET")
	api.HandleFunc("/user/{id}", controllers.GetUserById(db)).Methods("GET")
	api.HandleFunc("/users/{id}", controllers.UpdateUser(db)).Methods("PATCH")
	api.HandleFunc("/users/{id}", controllers.DeleteUser(db)).Methods("DELETE")
	// user := &models.User{ID: 1, FirstName: "man", LastName: "manega", Email: "man@man", Password: "manegaga"}
	// fmt.Println(user)
	log.Fatal(http.ListenAndServe(":8080", r))

}
