package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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
			break
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

	r := mux.NewRouter()

	appointmentRepo := repository.NewAppointmentRepository(db)
	appointmentHandler := handler.NewAppointmentHandler(appointmentRepo)

	medicalJustificationRepo := repository.NewMedicalJustificationRepository(db)
	medicalJustificationHandler := handler.NewMedicalJustificationHandler(medicalJustificationRepo)

	r.Handle("/getAllJustifications", http.HandlerFunc(medicalJustificationHandler.GetAllJustifications)).Methods("GET")
	r.Handle("/getJustificationsForParent/{userId}", http.HandlerFunc(medicalJustificationHandler.GetJustificationsForParent)).Methods("GET")
	r.Handle("/getJustificationsForDoctor/{userId}", http.HandlerFunc(medicalJustificationHandler.GetJustificationsForDoctor)).Methods("GET")
	r.Handle("/createJustification", http.HandlerFunc(medicalJustificationHandler.CreateJustification)).Methods("POST")

	// ---- Appointment routes ----
	r.Handle("/createAppointment", http.HandlerFunc(appointmentHandler.CreateAppointment)).Methods("POST")
	r.Handle("/getAppointmentsByParent/{id}", http.HandlerFunc(appointmentHandler.GetAppointmentsByParent)).Methods("GET")
	r.Handle("/getAppointmentsByDoctor/{id}", http.HandlerFunc(appointmentHandler.GetAppointmentsByDoctor)).Methods("GET")
	r.Handle("/getAppointments", http.HandlerFunc(appointmentHandler.GetAppointments)).Methods("GET")
	r.Handle("/getAppointmentsByDoctor", http.HandlerFunc(appointmentHandler.GetAppointmentsByDoctor)).Methods("GET")
	r.HandleFunc("/cancelAppointment/{id}", appointmentHandler.CancelAppointment).Methods("DELETE")
	r.HandleFunc("/justifyAppointment/{id}", appointmentHandler.JustifyAppointment).Methods("PUT")

	//r.Handle("/medicalRecord/user/{userId}", enableCORS(http.HandlerFunc(medicalJustificationHandler.GetJustificationsForParent))).Methods("GET")
	//r.Handle("/getJustification", enableCORS(http.HandlerFunc(medicalJustificationHandler.GetJustifications))).Methods("GET")
	//r.Handle("/createJustification", enableCORS(http.HandlerFunc(medicalJustificationHandler.CreateJustification))).Methods("POST")
	//
	//// ---- Appointment routes ----
	//r.Handle("/createAppointment", enableCORS(http.HandlerFunc(appointmentHandler.CreateAppointment))).Methods("POST")
	//r.Handle("/getAppointments/{id}", enableCORS(http.HandlerFunc(appointmentHandler.GetAppointments))).Methods("GET")
	//r.Handle("/getAppointments", enableCORS(http.HandlerFunc(appointmentHandler.GetAppointment))).Methods("GET")
	//r.Handle("/getAppointmentsByDoctor", enableCORS(http.HandlerFunc(appointmentHandler.GetAppointmentsByDoctor))).Methods("GET")

	//http.Handle("/getJustification", enableCORS(http.HandlerFunc(medicalJustificationHandler.GetJustifications)))
	//http.Handle("/medicalRecord/user/{userId}", enableCORS(http.HandlerFunc(medicalJustificationHandler.GetJustificationsForParent)))
	//http.Handle("/createJustification", enableCORS(http.HandlerFunc(medicalJustificationHandler.CreateJustification)))
	//
	//http.Handle("/createAppointment", enableCORS(http.HandlerFunc(appointmentHandler.CreateAppointment)))
	//http.Handle("/getAppointments/{id}", enableCORS(http.HandlerFunc(appointmentHandler.GetAppointments)))
	//http.Handle("/getAppointments", enableCORS(http.HandlerFunc(appointmentHandler.GetAppointment)))
	//http.Handle("/getAppointmentsByDoctor", enableCORS(http.HandlerFunc(appointmentHandler.GetAppointmentsByDoctor)))
	//
	//log.Println("Server running on :8081")
	//log.Fatal(http.ListenAndServe(":8081", nil))

	//http.ListenAndServe(":8081", r)
	log.Println("Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", enableCORS(r)))

}

//func enableCORS(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
//		w.Header().Set("Access-Control-Allow-Credentials", "true")
//
//		if r.Method == http.MethodOptions {
//			w.WriteHeader(http.StatusNoContent)
//			return
//		}
//		next.ServeHTTP(w, r)
//	})
//}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
