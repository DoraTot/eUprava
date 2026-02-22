package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"main.go/model"
	"main.go/repository"
	"net/http"
	"time"
)

type ChildrenHandler struct {
	Repo *repository.ChildRepo
}

func NewChildrenHandler(repo *repository.ChildRepo) *ChildrenHandler {
	return &ChildrenHandler{Repo: repo}
}

func (h *ChildrenHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	children, err := h.Repo.GetAllChildren()
	if err != nil {
		http.Error(w, "Failed to fetch children", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(children)
}

func (h *ChildrenHandler) GetByParentID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	child, err := h.Repo.GetChildByParentID(id)
	if err != nil {
		http.Error(w, "Child not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(child)
}

func (h *ChildrenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var a model.Child
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existing, err := h.Repo.GetChildByParentAndName(a.ParentId, a.Name)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, "Child with the same name already exists for this parent", http.StatusConflict)
		return
	}

	a.Enrolled = model.Pending
	id, err := h.Repo.Create(a.Name, a.ParentId, string(a.Enrolled))
	if err != nil {
		http.Error(w, "Failed to insert child: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     id,
	})
}

func (h *ChildrenHandler) EnrollChild(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParentID  string `json:"parent_id"`
		ChildName string `json:"child_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		fmt.Println("Decode error:", err)
		return
	}

	if req.ParentID == "" || req.ChildName == "" {
		http.Error(w, "parent_id and child_name are required", http.StatusBadRequest)
		fmt.Println("Missing parent_id or child_name")
		return
	}

	healthURL := "http://healthcare:8081"
	examDone := false

	type ExamCheckRequest struct {
		ParentID  string `json:"parent_id"`
		ChildName string `json:"child_name"`
	}
	type ExamCheckResponse struct {
		ParentID   string    `json:"parent_id"`
		ChildName  string    `json:"child_name"`
		Status     string    `json:"status"` // SCHEDULED, COMPLETED, EXPIRED
		ValidUntil time.Time `json:"valid_until"`
	}

	payload := ExamCheckRequest{
		ParentID:  req.ParentID,
		ChildName: req.ChildName,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal payload: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Marshal error:", err)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/exam/check", healthURL), bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Request creation error:", err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	fmt.Println("Sending request to Healthcare service:", httpReq.URL)

	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "Failed to contact Healthcare service: "+err.Error(), http.StatusServiceUnavailable)
		fmt.Println("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Healthcare response status:", resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Healthcare response body:", string(body))

	if resp.StatusCode == http.StatusOK {
		var examResp ExamCheckResponse
		if err := json.Unmarshal(body, &examResp); err != nil {
			http.Error(w, "Failed to decode Healthcare response: "+err.Error(), http.StatusInternalServerError)
			fmt.Println("Decode Healthcare response error:", err)
			return
		}

		fmt.Printf("Received exam info: %+v\n", examResp)

		if examResp.Status == "COMPLETED" && examResp.ValidUntil.After(time.Now()) {
			examDone = true
			_, _ = h.Repo.UpdateExamStatus(req.ParentID, req.ChildName, examDone)
			_, _ = h.Repo.UpdateEnrollmentStatus(req.ParentID, req.ChildName, string(model.Enrolled))
			fmt.Println("Child enrolled successfully")
		} else {
			fmt.Println("Exam not completed or expired")
		}
	} else {
		http.Error(w, "Healthcare service returned non-200 status", resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"exam_done": examDone,
		"parent_id": req.ParentID,
		"childName": req.ChildName,
	})
}

func (h *ChildrenHandler) ExamCompleted(w http.ResponseWriter, r *http.Request) {

	var req struct {
		ParentID  string `json:"parent_id"`
		ChildName string `json:"child_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ParentID == "" || req.ChildName == "" {
		http.Error(w, "parent_id and child_name are required", http.StatusBadRequest)
		return
	}

	_, err := h.Repo.UpdateEnrollmentStatus(req.ParentID, req.ChildName, string(model.Enrolled))
	if err != nil {
		http.Error(w, "Failed to update enrollment status", http.StatusInternalServerError)
		return
	}
	_, _ = h.Repo.UpdateExamStatus(req.ParentID, req.ChildName, true)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "enrolled",
		"parent_id":  req.ParentID,
		"child_name": req.ChildName,
	})
}

func (h *ChildrenHandler) ExamExpired(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParentID  string `json:"parent_id"`
		ChildName string `json:"child_name"`
		Status    string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("Exam expired for:", req.ChildName)

	_, err := h.Repo.UpdateEnrollmentStatus(req.ParentID, req.ChildName, string(model.NotEnrolled))
	if err != nil {
		http.Error(w, "Failed to update child status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "child status updated to exam expired",
	})
}

func (h *ChildrenHandler) GetChildsEnrollmentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	childName := vars["childName"]
	parentID := vars["parentId"]

	if childName == "" || parentID == "" {
		http.Error(w, "childName and parentId are required", http.StatusBadRequest)
		return
	}

	child, err := h.Repo.GetChildByParentAndName(parentID, childName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Child not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"childName": child.Name,
		"parentID":  child.ParentId,
		"status":    child.Enrolled,
	})
}
