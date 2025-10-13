package repository

import (
	"database/sql"
	"main.go/model"
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
