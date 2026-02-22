package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"log"
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

	createdJ, err := h.Repo.CreateJustification(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendanceURL := "http://app:8080/attendance/justify"
	payload := map[string]string{
		"parent":     createdJ.ParentID,
		"child":      createdJ.ChildName,
		"valid_from": createdJ.ValidFrom,
		"valid_to":   createdJ.ValidTo,
	}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(attendanceURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		log.Println("Failed to notify Attendance service:", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Println("Attendance service returned status:", resp.Status)
		}
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
	doctorID := vars["id"]
	log.Println("Received doctorID:", doctorID)

	justification, err := h.Repo.GetAppointmentsByDoctor(doctorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(justification)
}

func (h *MedicalJustificationHandler) DownloadJustificationPDF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	log.Println("Received med justification id:", idStr)

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	justification, err := h.Repo.GetJustificationByID(idInt)
	if err != nil || justification == nil {
		http.Error(w, "Justification not found", http.StatusNotFound)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Medical Justification", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 12)
	lines := []string{
		fmt.Sprintf("ID: %d", justification.ID),
		fmt.Sprintf("Child: %s", justification.ChildName),
		fmt.Sprintf("Doctor ID: %s", justification.DoctorID),
		fmt.Sprintf("Parent ID: %s", justification.ParentID),
		fmt.Sprintf("Valid From: %s", justification.ValidFrom),
		fmt.Sprintf("Valid To: %s", justification.ValidTo),
		fmt.Sprintf("Reason: %s", justification.Reason),
	}
	for _, line := range lines {
		pdf.CellFormat(0, 8, line, "", 1, "", false, 0, "")
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		log.Println("Failed to generate PDF:", err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="justification_%d.pdf"`, justification.ID))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Println("Failed to write PDF to response:", err)
	}
}
