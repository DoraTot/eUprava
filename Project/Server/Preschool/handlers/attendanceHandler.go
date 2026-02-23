package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main.go/repository"
	ws "main.go/websocket"
	"net/http"
	"strings"
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

func (h *AttendanceHandler) GetRecordsForToday(w http.ResponseWriter, r *http.Request) {
	records, err := h.Repo.GetAllAttendanceForToday()
	if err != nil {
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func (h *AttendanceHandler) GetRecordsByParent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentID := vars["id"]
	if parentID == "" {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}
	records, err := h.Repo.GetAllAttendanceByParent(parentID)
	if err != nil {
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func (h *AttendanceHandler) GetRecordsByParentForToday(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentID := vars["id"]
	if parentID == "" {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}
	records, err := h.Repo.GetAllAttendanceByParentForToday(parentID)
	if err != nil {
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func (h *AttendanceHandler) PickUp(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Parent    string `json:"parent"`
		ChildName string `json:"childName"`
		Date      string `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	att, err := h.Repo.PickUp(req.Parent, req.ChildName, req.Date, true)
	if err != nil {
		log.Println("Failed to insert attendance:", err)
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}

	message := att.Child + " was picked up from preschool at " + att.Date

	log.Println("Sending notification to parent:", att.Parent, "Message:", message)

	ws.SendNotification(att.Parent, message)
	log.Println("Notification sent successfully to parent:", att.Parent)
	notification := map[string]interface{}{
		"user_id": att.Parent,
		"message": message,
	}

	jsonData, _ := json.Marshal(notification)

	resp, err := http.Post("http://app:8080/notifications",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Failed to send notification to preschool service:", err)
	} else {
		defer resp.Body.Close()
		log.Println("Notification sent to preschool service, status:", resp.Status)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Picked up status updated",
	})

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
		Child     string `json:"child"`
		Parent    string `json:"parent"`
		Date      string `json:"date"`
		Missing   bool   `json:"missing"`
		PickedUp  bool   `json:"picked_up"`
		Justified bool   `json:"justified"`
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

	if req.Missing {
		healthcareURL := "http://healthcare:8081/checkJustificationForDate/" +
			req.Parent + "/" + req.Child

		resp, err := http.Get(healthcareURL)
		if err != nil {
			http.Error(w, "Failed to contact healthcare service", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			req.Justified = true
		} else if resp.StatusCode == http.StatusNotFound {
			req.Justified = false
		} else {
			http.Error(w, "Healthcare service error", http.StatusInternalServerError)
			return
		}
	}

	recordID, err := h.Repo.InsertAttendance(req.Child, req.Parent, dateTime, req.Missing, false, req.Justified)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		log.Println("Failed to insert attendance:", err)
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}

	att, err := h.Repo.GetByID(recordID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if att.Missing {
		message := ""
		if att.Justified {
			message = "Your child " + att.Child + " was absent from preschool today, their absence is justified with a valid medical justification."

		} else {
			message = "Your child " + att.Child + " was absent from preschool today, please provide medical justification."
		}

		log.Println("Sending notification to parent:", att.Parent, "Message:", message)

		ws.SendNotification(att.Parent, message)
		log.Println("Notification sent successfully to parent:", att.Parent)
		notification := map[string]interface{}{
			"user_id": att.Parent,
			"message": message,
		}

		jsonData, _ := json.Marshal(notification)

		resp, err := http.Post("http://localhost:8080/notifications",
			"application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Failed to send notification to preschool service:", err)
		} else {
			defer resp.Body.Close()
			log.Println("Notification sent to preschool service, status:", resp.Status)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Attendance added",
		"id":      recordID,
	})
}

func (h *AttendanceHandler) JustifyAttendance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Parent    string `json:"parent"`
		Child     string `json:"child"`
		ValidFrom string `json:"valid_from"`
		ValidTo   string `json:"valid_to"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	validFrom, err := time.Parse("2006-01-02", req.ValidFrom)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	validTo, err := time.Parse("2006-01-02", req.ValidTo)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	log.Println("med service sent:", req.Parent, req.Child, validFrom, validTo)

	rows, err := h.Repo.JustifyAttendance(req.Parent, req.Child, validFrom, validTo)
	if err != nil {
		log.Println("Failed to justify attendance:", err)
		http.Error(w, "Failed to justify attendance", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "No attendance record found for given parent and date", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Attendance justified",
		"updated": rows,
	})
}

func (h *AttendanceHandler) GetChildStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	childName := vars["childName"]
	if childName == "" {
		http.Error(w, "Missing childName", http.StatusBadRequest)
		return
	}

	stats, err := h.Repo.GetChildStatistics(childName)
	if err != nil {
		http.Error(w, "Failed to get attendance statistics", http.StatusInternalServerError)
		return
	}

	healthcareURL := "http://healthcare:8081/getChildStats/" + childName
	resp, err := http.Get(healthcareURL)
	if err != nil {
		http.Error(w, "Failed to contact healthcare service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Healthcare service returned an error", http.StatusInternalServerError)
		return
	}

	var healthcareStats struct {
		MonthlyAppointments int `json:"monthlyAppointments"`
		YearlyAppointments  int `json:"yearlyAppointments"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&healthcareStats); err != nil {
		http.Error(w, "Failed to decode healthcare stats", http.StatusInternalServerError)
		return
	}

	stats.MonthlyAppointments = healthcareStats.MonthlyAppointments
	stats.YearlyAppointments = healthcareStats.YearlyAppointments

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
