package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main.go/model"
	"main.go/repository"
	"net/http"
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

	err := h.Repo.DeleteAppointment(appointmentID)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "justified"})
}
