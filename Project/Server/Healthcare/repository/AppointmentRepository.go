package repository

import (
	"database/sql"
	"log"
	"main.go/model"
)

type AppointmentRepository struct {
	DB *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{DB: db}
}

func (r *AppointmentRepository) CreateAppointment(a *model.Appointment) error {
	query := `
		INSERT INTO appointments (child_name, parent_id, doctor_id, date_time, notes)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.DB.Exec(query, a.ChildName, a.ParentID, a.DoctorID, a.DateTime, a.Notes)
	return err
}

func (r *AppointmentRepository) GetAppointmentsByParent(parentID int) ([]model.Appointment, error) {
	query := `
		SELECT id, child_name, parent_id, doctor_id, date_time, notes
		FROM appointments
		WHERE parent_id = ?
	`
	rows, err := r.DB.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []model.Appointment
	for rows.Next() {
		var a model.Appointment
		if err := rows.Scan(&a.ID, &a.ChildName, &a.ParentID, &a.DoctorID, &a.DateTime, &a.Notes); err != nil {
			log.Println(err)
			continue
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

func (r *AppointmentRepository) GetAppointmentsByDoctor(doctorID int) ([]model.Appointment, error) {
	query := `
		SELECT id, child_name, parent_id, doctor_id, date_time, notes
		FROM appointments
		WHERE doctor_id = ?
	`
	rows, err := r.DB.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []model.Appointment
	for rows.Next() {
		var a model.Appointment
		if err := rows.Scan(&a.ID, &a.ChildName, &a.ParentID, &a.DoctorID, &a.DateTime, &a.Notes); err != nil {
			continue
		}
		appointments = append(appointments, a)
	}

	return appointments, nil
}
