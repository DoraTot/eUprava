package repository

import (
	"database/sql"
	"fmt"
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

func (r *AttendanceRepo) GetAllAttendanceForToday() ([]model.AttendanceRecord, error) {
	query := `
        SELECT id, child, parent_auth0_id, date, missing, justified, picked_up
        FROM attendance_record
        WHERE DATE(date) = CURRENT_DATE
        ORDER BY date DESC
    `
	rows, err := r.DB.Query(query)
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

func (r *AttendanceRepo) GetAllAttendanceByParentForToday(parentID string) ([]model.AttendanceRecord, error) {
	query := `
		SELECT id, child, parent_auth0_id, date, missing, justified, picked_up
		FROM attendance_record
		WHERE parent_auth0_id = ? AND DATE(date) = CURRENT_DATE
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
	var exists int
	checkQuery := `
        SELECT 1
        FROM attendance_record
        WHERE child = ? AND DATE(date) = DATE(?)
        LIMIT 1
    `
	err := r.DB.QueryRow(checkQuery, child, date).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking existing attendance for child '%s' on '%s': %v", child, date.Format("2006-01-02"), err)
		return 0, err
	}
	if exists == 1 {
		return 0, fmt.Errorf("attendance for child '%s' on %s already exists", child, date.Format("2006-01-02"))
	}

	query := `INSERT INTO attendance_record (child, parent_auth0_id, date, missing, picked_up, justified) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, child, parentAuth0ID, date, missing, pickedUp, justified)
	if err != nil {
		log.Printf("Failed to insert attendance for child '%s' on '%s': %v", child, date.Format("2006-01-02"), err)
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

func (r *AttendanceRepo) GetChildStatistics(childID string) (model.ChildStatistics, error) {
	var stats model.ChildStatistics

	query := "SELECT missing, justified FROM attendance_record WHERE child = ?"
	rows, err := r.DB.Query(query, childID)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var missing, justified bool
		if err := rows.Scan(&missing, &justified); err != nil {
			return stats, err
		}
		stats.TotalRecords++
		if missing {
			stats.TotalAbsent++
		} else {
			stats.TotalPresent++
		}
		if justified {
			stats.TotalJustified++
		}
	}

	return stats, nil
}
