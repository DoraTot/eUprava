package model

type UserType string

const (
	Patient UserType = "Patient"
	Doctor  UserType = "Doctor"
)

type User struct {
	Fname    string   `json:"fname" db:"fname"`
	Lname    string   `json:"lname" db:"lname"`
	Auth0ID  string   `json:"auth0_id" db:"auth0_id"`
	Usertype UserType `json:"usertype" db:"usertype"`
}

type AttendanceRecord struct {
	ID        int    `json:"id" db:"id"`
	Child     string `json:"child" db:"child"`
	Parent    string `json:"parent" db:"parent"`
	Date      string `json:"date" db:"date"`
	Missing   bool   `json:"missing" db:"missing"`
	Justified bool   `json:"justified" db:"justified"`
}
