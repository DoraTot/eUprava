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

	//userRepo := repository.NewUserRepo(db)
	//userHandler := handlers.NewUserHandler(userRepo)

	attendanceRepo := repository.NewAttendanceRepo(db)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceRepo)
	childrenRepo := repository.NewChildrenRepo(db)
	if err := childrenRepo.EnsureTableExists(); err != nil {
		log.Fatal("Cannot ensure children table:", err)
	}
	childrenHandler := handlers.NewChildrenHandler(childrenRepo)

	r := mux.NewRouter()

	//http.Handle("/parents", enableCORS(http.HandlerFunc(userHandler.GetParents)))
	//http.Handle("/attendance", enableCORS(http.HandlerFunc(attendanceHandler.GetRecords)))
	//http.Handle("/attendance", enableCORS(http.HandlerFunc(attendanceHandler.PostRecord)))
	//r.Handle("/attendance", enableCORS(http.HandlerFunc(attendanceHandler.GetRecords))).Methods("POST")

	//r.Handle("/attendance/pickUp", enableCORS(http.HandlerFunc(attendanceHandler.PickUp))).Methods("POST")
	//r.Handle("/attendance", enableCORS(http.HandlerFunc(attendanceHandler.GetRecords))).Methods(http.MethodGet)
	//r.Handle("/attendance", enableCORS(http.HandlerFunc(attendanceHandler.PostRecord))).Methods(http.MethodPost)

	r.HandleFunc("/attendance", attendanceHandler.GetRecords).Methods("GET")
	r.HandleFunc("/attendanceByParent/{id}", attendanceHandler.GetRecordsByParent).Methods("GET")
	r.HandleFunc("/attendance", attendanceHandler.PostRecord).Methods("POST")
	r.HandleFunc("/attendance/pickUp", attendanceHandler.PickUp).Methods("POST")
	r.HandleFunc("/attendance/justify", attendanceHandler.JustifyAttendance).Methods("POST")

	r.HandleFunc("/children/add", childrenHandler.Create).Methods("POST")
	r.HandleFunc("/enrollChild", childrenHandler.EnrollChild).Methods("POST")
	r.HandleFunc("/examCompleted", childrenHandler.ExamCompleted).Methods("POST")
	r.HandleFunc("/children/getByParentID/{id}", childrenHandler.GetByParentID).Methods("GET")
	r.HandleFunc("/children/getAll", childrenHandler.GetAll).Methods("GET")

	//http.Handle("/attendance", enableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	switch r.Method {
	//	case http.MethodGet:
	//		attendanceHandler.GetRecords(w, r)
	//	case http.MethodPost:
	//		attendanceHandler.PostRecord(w, r)
	//		//default:
	//		//	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//	}
	//})))

	//http.Handle("/login", enableCORS(http.HandlerFunc(userHandler.HandleAuth0Login)))

	//log.Println("Server running on :8080")
	//log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(r)))

}

//func enableCORS(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//		if r.Method == "OPTIONS" {
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
