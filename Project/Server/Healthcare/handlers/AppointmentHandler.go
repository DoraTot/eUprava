package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main.go/model"
	"main.go/repository"
	ws "main.go/websocket"
	"net/http"
	"strconv"
	"time"
)

type AppointmentHandler struct {
	Repo *repository.AppointmentRepository
}

func NewAppointmentHandler(repo *repository.AppointmentRepository) *AppointmentHandler {
	return &AppointmentHandler{Repo: repo}
}

func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var a model.Appointment
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if a.SysExam {
		resp, err := http.Get("http://app:8080/child/status/childName/" + a.ChildName + "/parentId/" + a.ParentID)
		if err != nil {
			http.Error(w, "Failed to validate child status", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Child not found", http.StatusBadRequest)
			return
		}

		var statusResponse struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
			http.Error(w, "Invalid response from preschool service", http.StatusInternalServerError)
			return
		}

		if statusResponse.Status == "ENROLLED" {
			http.Error(w, "Systematic exam not allowed for this child", http.StatusBadRequest)
			return
		}
	}

	createdA, err := h.Repo.CreateAppointment(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if a.SysExam {
		sysExam := model.SystematicExam{
			ChildName:     a.ChildName,
			ParentID:      a.ParentID,
			DoctorID:      a.DoctorID,
			Date:          time.Now(),
			ValidUntil:    nil,
			Status:        model.ExamScheduled,
			AppointmentID: createdA.ID,
		}

		jsonData, _ := json.Marshal(sysExam)
		resp, err := http.Post("http://localhost:8081/createSystematicExam",
			"application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Failed to call health service:", err)
			return
		}
		defer resp.Body.Close()
	}

	//w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *AppointmentHandler) GetAppointmentsByParent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentID := vars["id"]
	if parentID == "" {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}

	appointments, err := h.Repo.GetAppointmentsByParent(parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(appointments)
}

func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {

	appointments, err := h.Repo.GetAppointments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

func (h *AppointmentHandler) GetAppointmentsByDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorID := vars["id"]
	if doctorID == "" {
		http.Error(w, "Missing doctor ID", http.StatusBadRequest)
		return
	}

	appointments, err := h.Repo.GetAppointmentsByDoctor(doctorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(appointments)
}

func (h *AppointmentHandler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentID := vars["id"]
	if appointmentID == "" {
		http.Error(w, "Missing appointment ID", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(appointmentID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	app, err := h.Repo.GetAppointmentByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	message := "Parent of " + app.ChildName + " cancelled their appointment."

	log.Println("Sending notification to doctor:", app.DoctorID, "Message:", message)

	ws.SendNotification(app.DoctorID, message)
	log.Println("Notification sent successfully to doctor:", app.DoctorID)
	notification := map[string]interface{}{
		"user_id": app.DoctorID,
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

	err = h.Repo.DeleteAppointment(appointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "canceled"})
}

func (h *AppointmentHandler) JustifyAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentID := vars["id"]
	if appointmentID == "" {
		http.Error(w, "Missing appointment ID", http.StatusBadRequest)
		return
	}

	err := h.Repo.SetAppointmentJustified(appointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idInt, err := strconv.Atoi(appointmentID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	app, err := h.Repo.GetAppointmentByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	message := "The medical justification has been successfully delivered to the preschool your child " +
		app.ChildName +
		" attends. Their absence is justified for the duration of the medical justification validity period."
	log.Println("Sending notification to parent:", app.ParentID, "Message:", message)

	ws.SendNotification(app.ParentID, message)
	log.Println("Notification sent successfully to parent:", app.ParentID)
	notification := map[string]interface{}{
		"user_id": app.ParentID,
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
	json.NewEncoder(w).Encode(map[string]string{"status": "justified"})
}

func (h *AppointmentHandler) JustifyAppointmentWithoutNotif(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentID := vars["id"]
	if appointmentID == "" {
		http.Error(w, "Missing appointment ID", http.StatusBadRequest)
		return
	}

	err := h.Repo.SetAppointmentJustified(appointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "justified"})
}

func (h *AppointmentHandler) GetChildStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	childName := vars["childName"]
	if childName == "" {
		http.Error(w, "Missing childID", http.StatusBadRequest)
		return
	}

	stats, err := h.Repo.GetChildAppointmentStatistics(childName)
	if err != nil {
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
