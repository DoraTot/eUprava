package handler

import (
	"encoding/json"
	"main.go/model"
	"main.go/repository"
	"net/http"
	"strconv"
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

func (h *MedicalJustificationHandler) GetJustifications(w http.ResponseWriter, r *http.Request) {
	parentIDStr := r.URL.Query().Get("parentId")
	parentID, err := strconv.Atoi(parentIDStr)
	if err != nil {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}

	justifications, err := h.Repo.GetJustificationsByParent(parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(justifications)
}
