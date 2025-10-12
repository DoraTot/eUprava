package repository

import (
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"main.go/config"
	"main.go/model"
	"main.go/utils"
	"net/http"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(user *model.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (fname, lname, username, password, usertype)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = r.DB.Exec(query, user.Fname, user.Lname, user.Username, hashedPassword, user.Usertype)
	return err
}

func GetAuth0PublicKey(kid string) (*rsa.PublicKey, error) {
	resp, err := http.Get("https://" + config.Auth0Domain + "/.well-known/jwks.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var jwksData model.Jwks
	if err := json.Unmarshal(body, &jwksData); err != nil {
		return nil, err
	}

	for _, key := range jwksData.Keys {
		if key.Kid == kid {
			return jwt.ParseRSAPublicKeyFromPEM([]byte(convertX5CToPEM(key.X5c[0])))
		}
	}
	return nil, fmt.Errorf("public key not found for kid %s", kid)
}

func convertX5CToPEM(x5c string) string {
	return "-----BEGIN CERTIFICATE-----\n" + x5c + "\n-----END CERTIFICATE-----"
}
