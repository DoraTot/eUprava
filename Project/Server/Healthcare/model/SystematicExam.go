package model

import "time"

type ExamStatus string

const (
	ExamScheduled ExamStatus = "SCHEDULED"
	ExamCompleted ExamStatus = "COMPLETED"
	ExamExpired   ExamStatus = "EXPIRED"
)

type SystematicExam struct {
	ID            int        `json:"id" db:"id"`
	ChildName     string     `json:"child_name" db:"child_name"`
	ParentID      string     `json:"parent_id" db:"parent_id"`
	DoctorID      string     `json:"doctor_id" db:"doctor_id"`
	Date          time.Time  `json:"date" db:"date"`
	ValidUntil    *time.Time `json:"valid_until" db:"valid_until"`
	Status        ExamStatus `json:"status" db:"status"`
	AppointmentID int        `json:"appointment_id" db:"appointment_id"`
}
