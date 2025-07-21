// TODO add middlewares to protected logged in routes

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	itemController "test.com/events/controllers/itemController"
	"test.com/events/controllers/userController"
	"test.com/events/database"
	authMiddleware "test.com/events/middleware/auth"
)

func initEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (skipping)")
	}
}

func main() {
	initEnvVariables()
	err := database.Connect()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/create-account", userController.CreateAccount).Methods("POST")
	router.HandleFunc("/login", userController.Login).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware.AuthMiddleware)

	protected.HandleFunc("/items", itemController.GetItems).Methods("GET")
	protected.HandleFunc("/items", itemController.PostItem).Methods("POST")

	http.ListenAndServe(":8080", router)
}
