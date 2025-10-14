package model

type AppointmentStatus string

type Appointment struct {
	ID        int    `json:"id" db:"id"`
	ChildName string `json:"child_name" db:"child_name"`
	ParentID  int    `json:"parent_id" db:"parent_id"`
	DoctorID  int    `json:"doctor_id" db:"doctor_id"`
	DateTime  string `json:"date_time" db:"date_time"`
	Notes     string `json:"notes,omitempty" db:"notes"`
}
