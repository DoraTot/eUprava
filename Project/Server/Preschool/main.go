package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"main.go/config"
	"main.go/handlers"
	"main.go/repository"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Hello, World!")

	//dsn := "root:secret@tcp(db:3306)/e_uprava"
	dsn := config.GetDSN()
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Println("Waiting for DB (sql.Open)...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			break // ✅ DB is ready
		}

		log.Println("Waiting for DB (ping)...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	userRepo := repository.NewUserRepo(db)
	userHandler := handlers.NewUserHandler(userRepo)

	http.Handle("/register", enableCORS(http.HandlerFunc(userHandler.RegisterUser)))
	http.Handle("/login", enableCORS(http.HandlerFunc(userHandler.HandleAuth0Login)))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

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
