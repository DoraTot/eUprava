package model

type MedicalJustification struct {
	ID        int    `json:"id" db:"id"`
	ChildName string `json:"child_name" db:"child_name"`
	DoctorID  int    `json:"doctor_id" db:"doctor_id"`
	ParentID  int    `json:"parent_id" db:"parent_id"`
	Date      string `json:"date" db:"date"`
	Reason    string `json:"reason" db:"reason"`
}
