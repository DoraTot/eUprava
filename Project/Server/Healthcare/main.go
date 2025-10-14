package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lestrrat-go/jwx/v2/jwk"
	"log"
	"main.go/config"
	"main.go/handlers"
	"main.go/repository"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Hello, World!")

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

	appointmentRepo := repository.NewAppointmentRepository(db)
	appointmentHandler := handler.NewAppointmentHandler(appointmentRepo)

	medicalJustificationRepo := repository.NewMedicalJustificationRepository(db)
	medicalJustificationHandler := handler.NewMedicalJustificationHandler(medicalJustificationRepo)

	http.Handle("/getJustification", enableCORS(http.HandlerFunc(medicalJustificationHandler.GetJustifications)))
	http.Handle("/createJustification", enableCORS(http.HandlerFunc(medicalJustificationHandler.CreateJustification)))

	http.Handle("/createAppointment", enableCORS(http.HandlerFunc(appointmentHandler.CreateAppointment)))
	http.Handle("/getAppointments/{id}", enableCORS(http.HandlerFunc(appointmentHandler.GetAppointments)))
	http.Handle("/getAppointmentsByDoctor", enableCORS(http.HandlerFunc(appointmentHandler.GetAppointmentsByDoctor)))

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
