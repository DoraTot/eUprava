package handlers

import (
	"encoding/json"
	"main.go/repository"
	"net/http"
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

func (h *AttendanceHandler) HandleCreateAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var record struct {
		Child  string `json:"child"`
		Parent string `json:"parent_auth0_id"`
		Date   string `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO attendance_record (child, parent_auth0_id, date, missing, justified)
        VALUES (?, ?, ?, FALSE, FALSE)
    `
	_, err := h.Repo.DB.Exec(query, record.Child, record.Parent, record.Date)
	if err != nil {
		http.Error(w, "Failed to create record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
