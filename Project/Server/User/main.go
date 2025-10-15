package main

import (
	"fmt"
	_ "github.com/lestrrat-go/jwx/v2/jwk"
	"log"
	"main.g/handlers"
	"main.g/repository"

	"net/http"
)

func main() {
	fmt.Println("Hello, World!")

	userRepo := repository.NewUserRepo()
	userHandler := handlers.NewUserHandler(userRepo)

	http.Handle("/login", enableCORS(http.HandlerFunc(userHandler.HandleAuth0Login)))
	http.Handle("/parents", enableCORS(http.HandlerFunc(userHandler.GetParents)))
	http.Handle("/doctors", enableCORS(http.HandlerFunc(userHandler.GetDoctors)))

	log.Println("Server running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))

}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
