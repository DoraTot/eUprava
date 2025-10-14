package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"main.go/repository"
	"net/http"
)

type UserHandler struct {
	Repo *repository.UserRepo
}

func NewUserHandler(repo *repository.UserRepo) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) HandleAuth0Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid")
		}
		return repository.GetAuth0PublicKey(kid)
	})

	if err != nil || !token.Valid {
		log.Println("Token verification failed:", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	email, _ := claims["email"].(string)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User verified",
		"email":   email,
	})

}

func (h *UserHandler) GetParents(w http.ResponseWriter, r *http.Request) {
	parents, err := repository.GetParentsFromAuth0()
	if err != nil {
		log.Println("Error fetching parents:", err)
		http.Error(w, "Failed to fetch parents", http.StatusInternalServerError)
		return
	}

	log.Println("Fetched parents:", parents)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parents)
}
