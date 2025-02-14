package main

import (
	"BE/config"
	"BE/handlers"
	"log"
	"net/http"

	corsHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/register", handlers.Register).Methods("POST")

	protected := router.PathPrefix("").Subrouter()
	protected.Use(handlers.ValidateTokenMiddleware)
	protected.HandleFunc("/createArticle", handlers.CreateArticle).Methods("POST")
	protected.HandleFunc("/getArticle", handlers.GetArticles).Methods("GET")

	corsHandler := corsHandlers.CORS(
		corsHandlers.AllowedOrigins([]string{"*"}),
		corsHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		corsHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Println("Server berjalan di port 8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", corsHandler(router)))
}
