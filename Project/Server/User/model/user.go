package model

import "github.com/MicahParks/keyfunc"

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

//type AttendanceRecord struct {
//	ID        int    `json:"id" db:"id"`
//	Child     string `json:"child" db:"child"`
//	Parent    string `json:"parent" db:"parent"`
//	Date      string `json:"date" db:"date"`
//	Missing   bool   `json:"missing" db:"missing"`
//	Justified bool   `json:"justified" db:"justified"`
//}

type ParentUser struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type AuthMiddleware struct {
	JWKS     *keyfunc.JWKS
	Issuer   string
	Audience string
}

type TokenRequest struct {
	Token string `json:"token"`
}

type Jwks struct {
	Keys []struct {
		Kid string   `json:"kid"`
		Kty string   `json:"kty"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	} `json:"keys"`
}
