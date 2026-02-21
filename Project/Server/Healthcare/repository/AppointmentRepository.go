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

func (r *AppointmentRepository) CreateAppointment(a *model.Appointment) (*model.Appointment, error) {
	query := `
		INSERT INTO appointments (child_name, parent_id, doctor_id, date_time, notes, justified, sys_exam)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.DB.Exec(query, a.ChildName, a.ParentID, a.DoctorID, a.DateTime, a.Notes, a.Justified, a.SysExam)
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	a.ID = int(id)
	return a, nil
}

func (r *AppointmentRepository) GetAppointmentsByParent(parentID string) ([]model.Appointment, error) {
	query := `
		SELECT id, child_name, parent_id, doctor_id, date_time, notes, justified, sys_exam
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
		if err := rows.Scan(&a.ID, &a.ChildName, &a.ParentID, &a.DoctorID, &a.DateTime, &a.Notes, &a.Justified, &a.SysExam); err != nil {
			log.Println(err)
			continue
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

func (r *AppointmentRepository) GetAppointments() ([]model.Appointment, error) {
	query := `
		SELECT id, child_name, parent_id, doctor_id, date_time, notes, justified, sys_exam
		FROM appointments
		`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []model.Appointment
	for rows.Next() {
		var a model.Appointment
		if err := rows.Scan(&a.ID, &a.ChildName, &a.ParentID, &a.DoctorID, &a.DateTime, &a.Notes, &a.Justified, &a.SysExam); err != nil {
			log.Println(err)
			continue
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

func (r *AppointmentRepository) GetAppointmentsByDoctor(doctorID string) ([]model.Appointment, error) {
	query := `
		SELECT id, child_name, parent_id, doctor_id, date_time, notes, justified, sys_exam
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
		if err := rows.Scan(&a.ID, &a.ChildName, &a.ParentID, &a.DoctorID, &a.DateTime, &a.Notes, &a.Justified, &a.SysExam); err != nil {
			continue
		}
		appointments = append(appointments, a)
	}

	return appointments, nil
}

func (r *AppointmentRepository) DeleteAppointment(id string) error {
	query := `DELETE FROM appointments WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *AppointmentRepository) SetAppointmentJustified(id string) error {
	query := `UPDATE appointments SET justified = true WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *AppointmentRepository) GetAppointmentByID(id string) (*model.Appointment, error) {
	query := `
		SELECT id, child_name, parent_id, doctor_id, date_time, notes, justified, sys_exam
		FROM appointments
		WHERE id = ?
	`
	row := r.DB.QueryRow(query, id)

	var a model.Appointment
	err := row.Scan(&a.ID, &a.ChildName, &a.ParentID, &a.DoctorID, &a.DateTime, &a.Notes, &a.Justified, &a.SysExam)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}
