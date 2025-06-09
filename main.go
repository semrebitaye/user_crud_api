package main

import (
	"log"
	"net/http"
	initializer "user_crud_api/initializers"

	"github.com/gorilla/mux"
)

func main() {
	initializer.LoadEnvVariable()
	initializer.Connect()

	r := mux.NewRouter()

	// user := &models.User{ID: 1, FirstName: "man", LastName: "manega", Email: "man@man", Password: "manegaga"}
	// fmt.Println(user)
	log.Fatal(http.ListenAndServe(":8080", r))

}
