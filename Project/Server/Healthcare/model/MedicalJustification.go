package model

type MedicalJustification struct {
	ID        int    `json:"id" db:"id"`
	ChildName string `json:"child_name" db:"child_name"`
	DoctorID  string `json:"doctor_id" db:"doctor_id"`
	ParentID  string `json:"parent_id" db:"parent_id"`
	ValidFrom string `json:"valid_from" db:"valid_from"`
	ValidTo   string `json:"valid_to" db:"valid_to"`
	Reason    string `json:"reason" db:"reason"`
}
