package model

type UserType string

const (
	Doctor  UserType = "Doctor"
	Parent  UserType = "Parent"
	Teacher UserType = "Teacher"
)

type User struct {
	Auth0ID  string   `json:"auth0_id" db:"auth0_id"`
	Fname    string   `json:"fname" db:"fname"`
	Lname    string   `json:"lname" db:"lname"`
	Usertype UserType `json:"usertype" db:"usertype"`
}
