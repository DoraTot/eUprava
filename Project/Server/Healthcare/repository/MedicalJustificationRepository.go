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

func (r *MedicalJustificationRepository) CreateJustification(j *model.MedicalJustification) (*model.MedicalJustification, error) {
	query := `
		INSERT INTO medical_justifications (child_name, doctor_id, parent_id, valid_from, valid_to, reason)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := r.DB.Exec(query, j.ChildName, j.DoctorID, j.ParentID, j.ValidFrom, j.ValidTo, j.Reason)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	j.ID = int(id)
	return j, nil
}

func (r *MedicalJustificationRepository) GetJustificationsByParent(parentID string) ([]model.MedicalJustification, error) {
	query := `
		SELECT id, child_name, doctor_id, parent_id, valid_from, valid_to, reason
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
		if err := rows.Scan(&j.ID, &j.ChildName, &j.DoctorID, &j.ParentID, &j.ValidFrom, &j.ValidTo, &j.Reason); err != nil {
			log.Println(err)
			continue
		}
		justifications = append(justifications, j)
	}
	return justifications, nil
}

func (r *MedicalJustificationRepository) GetAppointmentsByDoctor(doctorID string) ([]model.MedicalJustification, error) {
	query := `
		SELECT id, child_name, doctor_id, parent_id, valid_from, valid_to, reason
		FROM medical_justifications
		WHERE doctor_id = ?
	`
	rows, err := r.DB.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []model.MedicalJustification
	for rows.Next() {
		var a model.MedicalJustification
		if err := rows.Scan(&a.ID, &a.ChildName, &a.DoctorID, &a.ParentID, &a.ValidFrom, &a.ValidTo, &a.Reason); err != nil {
			log.Println(err)
			continue
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

func (r *MedicalJustificationRepository) GetAllJustifications() ([]model.MedicalJustification, error) {
	query := `
		SELECT id, child_name, doctor_id, parent_id, valid_from, valid_to, reason
		FROM medical_justifications
		
		`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var justifications []model.MedicalJustification
	for rows.Next() {
		var a model.MedicalJustification
		if err := rows.Scan(&a.ID, &a.ChildName, &a.DoctorID, &a.ParentID, &a.ValidFrom, &a.ValidTo, &a.Reason); err != nil {
			log.Println(err)
			continue
		}
		justifications = append(justifications, a)
	}
	return justifications, nil
}
