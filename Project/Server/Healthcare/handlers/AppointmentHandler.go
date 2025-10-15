package handler

import (
	"encoding/json"
	"main.go/model"
	"main.go/repository"
	"net/http"
	"strconv"
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

	if err := h.Repo.CreateAppointment(&a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	parentIDStr := r.URL.Query().Get("parentId")
	parentID, err := strconv.Atoi(parentIDStr)
	if err != nil {
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

func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {

	appointments, err := h.Repo.GetAppointments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

func (h *AppointmentHandler) GetAppointmentsByDoctor(w http.ResponseWriter, r *http.Request) {
	doctorIDStr := r.URL.Query().Get("doctorId")
	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	appointments, err := h.Repo.GetAppointmentsByDoctor(doctorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(appointments)
}
