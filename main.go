package main

import (
	initializer "user_crud_api/initializers"
)

func main() {
	initializer.LoadEnvVariable()
	initializer.Connect()

	// user := &models.User{ID: 1, FirstName: "man", LastName: "manega", Email: "man@man", Password: "manegaga"}
	// fmt.Println(user)

}
