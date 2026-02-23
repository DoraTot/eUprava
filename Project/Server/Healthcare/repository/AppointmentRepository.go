package repository

import (
	"database/sql"
	"log"
	"main.go/model"
	"time"
)

type AppointmentRepository struct {
	DB *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{DB: db}
}

func (r *AppointmentRepository) EnsureTableExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS appointments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		child_name VARCHAR(255) NOT NULL,
		parent_id VARCHAR(255) NOT NULL,
		doctor_id VARCHAR(255) NOT NULL,
		date_time DATETIME NOT NULL,
		notes VARCHAR(512),
		justified BOOLEAN NOT NULL DEFAULT FALSE,
		sys_exam BOOLEAN NOT NULL DEFAULT FALSE
	);
	`

	_, err := r.DB.Exec(query)
	if err != nil {
		log.Println("Failed to create appointments table:", err)
	}

	return err
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

func (r *AppointmentRepository) GetAppointmentByID(id int) (*model.Appointment, error) {
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

func (r *AppointmentRepository) GetChildAppointmentStatistics(childName string) (model.ChildStatistics, error) {
	var stats model.ChildStatistics

	log.Printf("Getting appointment stats for child '%s'", childName)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	log.Printf("Month start: %v, Year start: %v", monthStart, yearStart)

	query := `
        SELECT date_time
        FROM appointments
        WHERE child_name = ?
    `
	rows, err := r.DB.Query(query, childName)
	if err != nil {
		log.Printf("Failed to query appointments for child '%s': %v", childName, err)
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var dateStr string
		if err := rows.Scan(&dateStr); err != nil {
			log.Printf("Error scanning row for child '%s': %v", childName, err)
			return stats, err
		}

		date, err := time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			log.Printf("Error parsing date '%s' for child '%s': %v", dateStr, childName, err)
			return stats, err
		}

		if date.After(monthStart) || date.Equal(monthStart) {
			stats.MonthlyAppointments++
		}

		if date.After(yearStart) || date.Equal(yearStart) {
			stats.YearlyAppointments++
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error for child '%s': %v", childName, err)
		return stats, err
	}

	log.Printf("Statistics for child '%s': %+v", childName, stats)
	return stats, nil
}
