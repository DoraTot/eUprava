package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main.go/model"
	"main.go/repository"
	"net/http"
)

type MedicalJustificationHandler struct {
	Repo *repository.MedicalJustificationRepository
}

func NewMedicalJustificationHandler(repo *repository.MedicalJustificationRepository) *MedicalJustificationHandler {
	return &MedicalJustificationHandler{Repo: repo}
}

func (h *MedicalJustificationHandler) CreateJustification(w http.ResponseWriter, r *http.Request) {
	var j model.MedicalJustification
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Repo.CreateJustification(&j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *MedicalJustificationHandler) GetAllJustifications(w http.ResponseWriter, r *http.Request) {
	justifications, err := h.Repo.GetAllJustifications()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(justifications)
}

func (h *MedicalJustificationHandler) GetJustificationsForParent(w http.ResponseWriter, r *http.Request) {
	//parentIDStr := r.URL.Query().Get("userId")
	vars := mux.Vars(r)
	parentID := vars["userId"]
	log.Println("Received parentID:", parentID)

	justification, err := h.Repo.GetJustificationsByParent(parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(justification)
}

func (h *MedicalJustificationHandler) GetJustificationsForDoctor(w http.ResponseWriter, r *http.Request) {
	//parentIDStr := r.URL.Query().Get("userId")
	vars := mux.Vars(r)
	doctorID := vars["userId"]
	log.Println("Received doctorID:", doctorID)

	justification, err := h.Repo.GetAppointmentsByDoctor(doctorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(justification)
}
