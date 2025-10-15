package repository

import (
	"database/sql"
	"log"
	"main.go/model"
)

type MedicalJustificationRepository struct {
	DB *sql.DB
}

func NewMedicalJustificationRepository(db *sql.DB) *MedicalJustificationRepository {
	return &MedicalJustificationRepository{DB: db}
}

func (r *MedicalJustificationRepository) CreateJustification(j *model.MedicalJustification) error {
	query := `
		INSERT INTO medical_justifications (child_name, doctor_id, parent_id, dated, reason)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.DB.Exec(query, j.ChildName, j.DoctorID, j.ParentID, j.Date, j.Reason)
	return err
}

func (r *MedicalJustificationRepository) GetJustificationsByParent(parentID int) ([]model.MedicalJustification, error) {
	query := `
		SELECT id, child_name, doctor_id, parent_id, dated, reason
		FROM medical_justifications
		WHERE parent_id = ?
	`
	rows, err := r.DB.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var justifications []model.MedicalJustification
	for rows.Next() {
		var j model.MedicalJustification
		if err := rows.Scan(&j.ID, &j.ChildName, &j.DoctorID, &j.ParentID, &j.Date, &j.Reason); err != nil {
			log.Println(err)
			continue
		}
		justifications = append(justifications, j)
	}
	return justifications, nil
}

func (r *MedicalJustificationRepository) GetAppointmentsByParent(parentID string) ([]model.MedicalJustification, error) {
	query := `
		SELECT id, child_name, doctor_id, parent_id, dated, reason
		FROM medical_justifications
		WHERE parent_id = ?
	`
	rows, err := r.DB.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []model.MedicalJustification
	for rows.Next() {
		var a model.MedicalJustification
		if err := rows.Scan(&a.ID, &a.ChildName, &a.ParentID, &a.DoctorID, &a.Date, &a.Reason); err != nil {
			log.Println(err)
			continue
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}
