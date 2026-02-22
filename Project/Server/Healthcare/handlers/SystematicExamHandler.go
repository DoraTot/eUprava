package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"main.go/model"
	"main.go/repository"
	"net/http"
	"strconv"
	"time"
)

type SystematicExamHandler struct {
	Repo *repository.SystematicExamRepository
}

func NewSystematicExamHandler(repo *repository.SystematicExamRepository) *SystematicExamHandler {
	return &SystematicExamHandler{Repo: repo}
}

func (h *SystematicExamHandler) CheckExamStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParentID  string `json:"parent_id"`
		ChildName string `json:"child_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	fmt.Println("this endpoint is hit with " + req.ChildName + "and " + req.ParentID)

	exam, err := h.Repo.GetExamByChild(req.ParentID, req.ChildName)
	if err != nil {
		http.Error(w, "Exam not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"parent_id":   exam.ParentID,
		"child_name":  exam.ChildName,
		"status":      exam.Status,
		"valid_until": exam.ValidUntil,
	})
}

func (h *SystematicExamHandler) CreateExam(w http.ResponseWriter, r *http.Request) {
	var exam model.SystematicExam
	if err := json.NewDecoder(r.Body).Decode(&exam); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if exam.Status == "" {
		exam.Status = model.ExamScheduled
	}

	if err := h.Repo.CreateExam(&exam); err != nil {
		http.Error(w, "Failed to create exam: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *SystematicExamHandler) UpdateExamStatusByAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentIDStr := vars["appointmentId"]

	appointmentID, err := strconv.Atoi(appointmentIDStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newStatus := model.ExamStatus(req.Status)

	if err := h.Repo.UpdateExamStatusByAppointmentID(appointmentID, newStatus); err != nil {
		http.Error(w, "Failed to update exam status", http.StatusInternalServerError)
		return
	}

	if newStatus == model.ExamCompleted {
		go h.notifyPreschool(appointmentID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *SystematicExamHandler) notifyPreschool(appointmentID int) {

	exam, err := h.Repo.GetExamByAppointmentID(appointmentID)
	if err != nil {
		return
	}

	type PreschoolNotification struct {
		ParentID  string `json:"parent_id"`
		ChildName string `json:"child_name"`
	}

	payload := PreschoolNotification{
		ParentID:  exam.ParentID,
		ChildName: exam.ChildName,
	}

	data, _ := json.Marshal(payload)

	preschoolURL := "http://app:8080/examCompleted"

	req, _ := http.NewRequest("POST", preschoolURL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("Failed to notify preschool:", err)
	}
}

func (h *SystematicExamHandler) NotifyPreschoolExpired(exam *model.SystematicExam) {

	type PreschoolExpiredNotification struct {
		ParentID  string `json:"parent_id"`
		ChildName string `json:"child_name"`
		Status    string `json:"status"`
	}

	payload := PreschoolExpiredNotification{
		ParentID:  exam.ParentID,
		ChildName: exam.ChildName,
		Status:    string(model.ExamExpired),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal expired notification:", err)
		return
	}

	preschoolURL := "http://app:8080/examExpired"

	req, err := http.NewRequest("POST", preschoolURL, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to notify preschool (expired):", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Preschool notified about expired exam for:", exam.ChildName)
}

func (h *SystematicExamHandler) RunExpiredExamCheck(w http.ResponseWriter, r *http.Request) {
	exams, err := h.Repo.GetExpiredExams()
	if err != nil {
		http.Error(w, "Error fetching expired exams: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var expiredExams []map[string]interface{}

	for _, exam := range exams {
		if err := h.Repo.UpdateExamStatus(exam.ID, model.ExamExpired); err != nil {
			fmt.Println("Failed to update exam status:", err)
			continue
		}

		go h.NotifyPreschoolExpired(&exam)

		expiredExams = append(expiredExams, map[string]interface{}{
			"id":          exam.ID,
			"child_name":  exam.ChildName,
			"parent_id":   exam.ParentID,
			"valid_until": exam.ValidUntil,
			"status":      exam.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expiredExams)
}

func (h *SystematicExamHandler) CheckExpiredExams() {
	exams, err := h.Repo.GetExpiredExams()
	if err != nil {
		fmt.Println("Error fetching expired exams:", err)
		return
	}

	for _, exam := range exams {
		fmt.Println("Expiring exam for:", exam.ChildName)

		if err := h.Repo.UpdateExamStatus(exam.ID, model.ExamExpired); err != nil {
			fmt.Println("Failed to update exam status:", err)
			continue
		}

		go h.NotifyPreschoolExpired(&exam)
	}
}
