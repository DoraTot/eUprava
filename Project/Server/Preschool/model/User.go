package model

type User struct {
	//ID       int    `json:"id" db:"id"`
	Fname    string `json:"fname" db:"fname"`
	Lname    string `json:"lname" db:"lname"`
	Username string `json:"username" db:"username"`
	Password string `json:"password,omitempty" db:"password"`
	Usertype string `json:"usertype" db:"usertype"`
}
