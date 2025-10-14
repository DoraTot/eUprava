package repository

import (
	"database/sql"
	"main.go/model"
	"time"
)

type AttendanceRepo struct {
	DB *sql.DB
}

func NewAttendanceRepo(db *sql.DB) *AttendanceRepo {
	return &AttendanceRepo{DB: db}
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
		err := rows.Scan(&rec.ID, &rec.Child, &rec.Parent, &rec.Date, &rec.Missing, &rec.Justified)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *AttendanceRepo) InsertAttendance(child string, parentAuth0ID string, date time.Time, missing bool) (int64, error) {
	query := `INSERT INTO attendance_record (child, parent_auth0_id, date, missing) VALUES (?, ?, ?, ?)`
	res, err := r.DB.Exec(query, child, parentAuth0ID, date, missing)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
