package repository

import (
	"database/sql"
	"main.go/model"
	"main.go/utils"
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
