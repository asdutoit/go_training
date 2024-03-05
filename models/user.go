package models

import (
	"errors"

	"github.com/asdutoit/gotraining/section11/db"
	"github.com/asdutoit/gotraining/section11/utils"
)

type User struct {
	ID       int64
	Username string
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `
	INSERT INTO users(username, email, password) 
	VALUES (?, ?, ?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Username, u.Email, hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func GetUserByEmail(email string) (*User, error) {
	query := `
	SELECT id, username, email, password
	FROM users
	WHERE email = ?`

	row := db.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func GetAllUsers() ([]User, error) {
	query := `
	SELECT id, username, email
	FROM users`

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *User) ValidateCredentials() (int64, error) {
	user, err := GetUserByEmail(u.Email)

	if err != nil {
		return 0, err
	}

	valid := utils.CheckPasswordHash(u.Password, user.Password)

	if !valid {
		return 0, errors.New("invalid credentials")
	}

	return user.ID, nil
}
