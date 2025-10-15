package model

type MedicalJustification struct {
	ID        int    `json:"id" db:"id"`
	ChildName string `json:"child_name" db:"child_name"`
	DoctorID  string `json:"doctor_id" db:"doctor_id"`
	ParentID  string `json:"parent_id" db:"parent_id"`
	Date      string `json:"dated" db:"dated"`
	Reason    string `json:"reason" db:"reason"`
}
