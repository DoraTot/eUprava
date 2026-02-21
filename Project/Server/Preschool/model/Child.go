package model

type EnrolledStatus string

const (
	Enrolled    EnrolledStatus = "ENROLLED"
	Pending     EnrolledStatus = "PENDING"
	NotEnrolled EnrolledStatus = "NOTENROLLED"
)

type Child struct {
	ID       int            `json:"id" db:"id"`
	Name     string         `json:"name" db:"name"`
	ParentId string         `json:"parent_id" db:"parent_id"`
	ExamDone bool           `json:"exam_done" db:"exam_done"`
	Enrolled EnrolledStatus `json:"enrolled" db:"enrolled"`
}
