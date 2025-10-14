package handlers

import (
	"encoding/json"
	"log"
	"main.go/repository"
	"net/http"
	"time"
)

type AttendanceHandler struct {
	Repo *repository.AttendanceRepo
}

func NewAttendanceHandler(repo *repository.AttendanceRepo) *AttendanceHandler {
	return &AttendanceHandler{Repo: repo}
}

func (h *AttendanceHandler) GetRecords(w http.ResponseWriter, r *http.Request) {
	records, err := h.Repo.GetAllAttendance()
	if err != nil {
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func (h *AttendanceHandler) PostRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Child   string `json:"child"`
		Parent  string `json:"parent"`
		Date    string `json:"date"`
		Missing bool   `json:"missing"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	dateTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	recordID, err := h.Repo.InsertAttendance(req.Child, req.Parent, dateTime, req.Missing)
	if err != nil {
		log.Println("Failed to insert attendance:", err)
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Attendance added",
		"id":      recordID,
	})
}
