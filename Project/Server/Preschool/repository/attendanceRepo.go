package repository

import (
	"database/sql"
	"log"
	"main.go/model"
	"time"
)

type AttendanceRepo struct {
	DB *sql.DB
}

func NewAttendanceRepo(db *sql.DB) *AttendanceRepo {
	return &AttendanceRepo{DB: db}
}

func (r *AttendanceRepo) EnsureTableExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS attendance_record (
		id INT AUTO_INCREMENT PRIMARY KEY,
		child VARCHAR(255) NOT NULL,
		parent_auth0_id VARCHAR(255) NOT NULL,
		date DATETIME NOT NULL,
		missing BOOLEAN NOT NULL DEFAULT FALSE,
		justified BOOLEAN NOT NULL DEFAULT FALSE,
		picked_up BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err := r.DB.Exec(query)
	if err != nil {
		log.Println("Failed to create attendance_record table:", err)
	}

	return err
}

func (r *AttendanceRepo) GetAllAttendance() ([]model.AttendanceRecord, error) {
	rows, err := r.DB.Query("SELECT * FROM attendance_record ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.AttendanceRecord
	for rows.Next() {
		var rec model.AttendanceRecord
		err := rows.Scan(&rec.ID, &rec.Child, &rec.Parent, &rec.Date, &rec.Missing, &rec.Justified, &rec.PickedUp)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *AttendanceRepo) GetAllAttendanceByParent(parentID string) ([]model.AttendanceRecord, error) {
	query := `
		SELECT id, child, parent_auth0_id, date, missing, justified, picked_up
		FROM attendance_record
		WHERE parent_auth0_id = ?
		ORDER BY date DESC
	`

	rows, err := r.DB.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.AttendanceRecord
	for rows.Next() {
		var rec model.AttendanceRecord
		if err := rows.Scan(
			&rec.ID,
			&rec.Child,
			&rec.Parent,
			&rec.Date,
			&rec.Missing,
			&rec.Justified,
			&rec.PickedUp,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r *AttendanceRepo) InsertAttendance(child string, parentAuth0ID string, date time.Time, missing bool, pickedUp bool, justified bool) (int64, error) {
	query := `INSERT INTO attendance_record (child, parent_auth0_id, date, missing, picked_up, justified) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, child, parentAuth0ID, date, missing, pickedUp, justified)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AttendanceRepo) PickUp(parentAuth0ID string, date time.Time, pickedUp bool) (int64, error) {
	query := `
		UPDATE attendance_record
		SET picked_up = ?
		WHERE parent_auth0_id = ? AND date = ?
	`

	res, err := r.DB.Exec(query, pickedUp, parentAuth0ID, date)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *AttendanceRepo) JustifyAttendance(parentID, childID string, validFrom, validTo time.Time) (int64, error) {
	query := `
		UPDATE attendance_record
		SET justified = true
		WHERE parent_auth0_id = ? AND child = ? AND missing = true
		  AND date BETWEEN ? AND ?
	`
	res, err := r.DB.Exec(query, parentID, childID, validFrom, validTo)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
