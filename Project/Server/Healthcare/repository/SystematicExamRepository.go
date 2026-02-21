package repository

import (
	"database/sql"
	"fmt"
	"log"
	"main.go/model"
	"time"
)

type SystematicExamRepository struct {
	DB *sql.DB
}

func NewSystematicExamRepository(db *sql.DB) *SystematicExamRepository {
	return &SystematicExamRepository{DB: db}
}

func (r *SystematicExamRepository) EnsureTableExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS systematic_exams (
		id INT AUTO_INCREMENT PRIMARY KEY,
		child_name VARCHAR(255) NOT NULL,
		parent_id VARCHAR(255) NOT NULL,
		doctor_id VARCHAR(255) NOT NULL,
		date DATETIME NOT NULL,
		valid_until DATETIME NULL,
		status VARCHAR(50) NOT NULL,
		appointment_id INT NOT NULL
	);`
	_, err := r.DB.Exec(query)
	if err != nil {
		log.Println("Failed to create systematic_exams table:", err)
	}
	return err
}

func (r *SystematicExamRepository) CreateExam(exam *model.SystematicExam) error {
	query := `INSERT INTO systematic_exams 
	(child_name, parent_id, doctor_id, date, valid_until, status, appointment_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, exam.ChildName, exam.ParentID, exam.DoctorID, exam.Date, exam.ValidUntil, string(exam.Status), exam.AppointmentID)
	return err
}

func (r *SystematicExamRepository) GetExamByChild(parentID, childName string) (*model.SystematicExam, error) {
	fmt.Println("DEBUG: Looking for exam with parentID=", parentID, " childName=", childName)

	query := `SELECT id, child_name, parent_id, doctor_id, date, valid_until, status, appointment_id 
	          FROM systematic_exams 
			  WHERE parent_id = ? AND child_name = ? 
			  ORDER BY date DESC LIMIT 1`
	row := r.DB.QueryRow(query, parentID, childName)

	var exam model.SystematicExam
	var dateStr, validUntilStr sql.NullString
	var status string

	err := row.Scan(
		&exam.ID,
		&exam.ChildName,
		&exam.ParentID,
		&exam.DoctorID,
		&dateStr,
		&validUntilStr,
		&status,
		&exam.AppointmentID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("DEBUG: No exam found for", parentID, childName)
		} else {
			fmt.Println("DEBUG: Scan error:", err)
		}
		return nil, err
	}

	if dateStr.Valid {
		t, parseErr := time.Parse("2006-01-02 15:04:05", dateStr.String)
		if parseErr != nil {
			fmt.Println("DEBUG: Failed to parse date:", parseErr)
			return nil, parseErr
		}
		exam.Date = t
	}

	if validUntilStr.Valid {
		t, parseErr := time.Parse("2006-01-02 15:04:05", validUntilStr.String)
		if parseErr != nil {
			fmt.Println("DEBUG: Failed to parse valid_until:", parseErr)
			return nil, parseErr
		}
		exam.ValidUntil = &t
	} else {
		exam.ValidUntil = nil
	}

	exam.Status = model.ExamStatus(status)

	fmt.Println("DEBUG: Found exam:", exam)
	return &exam, nil
}

func (r *SystematicExamRepository) UpdateExamStatus(id int, status model.ExamStatus) error {
	query := `UPDATE systematic_exams SET status = ? WHERE id = ?`
	_, err := r.DB.Exec(query, string(status), id)
	return err
}

func (r *SystematicExamRepository) UpdateExamStatusByAppointmentID(appointmentID int, status model.ExamStatus) error {
	validUntil := time.Now().AddDate(1, 0, 0)
	query := `
		UPDATE systematic_exams
		SET status = ?, valid_until = ?
		WHERE appointment_id = ?
	`

	result, err := r.DB.Exec(query, status, validUntil, appointmentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SystematicExamRepository) GetExamByAppointmentID(appointmentID int) (*model.SystematicExam, error) {
	query := `
		SELECT id, child_name, parent_id, doctor_id, date, valid_until, status, appointment_id
		FROM systematic_exams
		WHERE appointment_id = ?
	`

	row := r.DB.QueryRow(query, appointmentID)

	var exam model.SystematicExam
	var dateStr, validUntilStr sql.NullString
	var status string

	err := row.Scan(
		&exam.ID,
		&exam.ChildName,
		&exam.ParentID,
		&exam.DoctorID,
		&dateStr,
		&validUntilStr,
		&status,
		&exam.AppointmentID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("DEBUG: No exam found for")
		} else {
			fmt.Println("DEBUG: Scan error:", err)
		}
		return nil, err
	}

	if dateStr.Valid {
		t, parseErr := time.Parse("2006-01-02 15:04:05", dateStr.String)
		if parseErr != nil {
			fmt.Println("DEBUG: Failed to parse date:", parseErr)
			return nil, parseErr
		}
		exam.Date = t
	}

	if validUntilStr.Valid {
		t, parseErr := time.Parse("2006-01-02 15:04:05", validUntilStr.String)
		if parseErr != nil {
			fmt.Println("DEBUG: Failed to parse valid_until:", parseErr)
			return nil, parseErr
		}
		exam.ValidUntil = &t
	} else {
		exam.ValidUntil = nil
	}

	exam.Status = model.ExamStatus(status)

	fmt.Println("DEBUG: Found exam:", exam)
	return &exam, nil
}
