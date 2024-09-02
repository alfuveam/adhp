package models

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type User struct {
	Id           string
	Completename string    `json:"completename"`
	Cpf          string    `json:"cpf"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Country      string    `json:"country"`
	CreateAt     time.Time `json:"createat"`
	LastLogin    time.Time `json:"lastlogin"`
	UserType     int       `json:"usertype"`
}

// func (u *User) VerifyEmail(db *sql.DB) {
// 	// conexao com banco de dados
// }

// func (u *User) CheckDuplicateEmail(db *sql.DB) (bool, error) {
// 	// conexao com banco de dados
// 	var enough bool
// 	// Query for a value based on a single row.
// 	if err := db.QueryRow("SELECT * from users where email = ?",
// 		u.Email).Scan(&enough); err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, fmt.Errorf("Email in use: %d:", u.Email)
// 		}
// 		return false, fmt.Errorf("Email in use: %d: %v", u.Email, err)
// 	}
// 	return enough, nil
// }

func (u *User) CreateUser(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (completename ,cpf ,phone ,email ,password ,country ,createat, lastlogin) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		u.Completename,
		u.Cpf,
		u.Phone,
		u.Email,
		u.Password,
		u.Country,
		u.CreateAt,
		u.LastLogin,
	)
	return err
}

func (u *User) OnLoginCheckEmailAndPasswordId(db *sql.DB) (string, bool, error) {
	var userID string
	err := db.QueryRow("SELECT id FROM users WHERE email = $1 AND password = $2", u.Email, u.Password).Scan(&userID)
	if err != nil {
		log.Println("OnLoginCheckEmailAndPassword: ", err)
		return "", false, errors.New("Verify email and password.")
	}
	return userID, true, nil
}

func (u *User) OnLoginCheckEmailAndPasswordUser(db *sql.DB) (User, bool, error) {
	var user User
	err := db.QueryRow("SELECT id, completename ,cpf ,phone ,email ,password ,country ,createat, lastlogin, usertype FROM users WHERE email = $1", u.Email).
		Scan(
			&user.Id,
			&user.Completename,
			&user.Cpf,
			&user.Phone,
			&user.Email,
			&user.Password,
			&user.Country,
			&user.CreateAt,
			&user.LastLogin,
			&user.UserType,
		)
	if err != nil {
		log.Println("OnLoginCheckEmailAndPassword: ", err)
		return user, false, errors.New("Verify email and password.")
	}
	return user, true, nil
}

func (u *User) GetUserById(db *sql.DB) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, completename ,cpf ,phone ,email ,password ,country ,createat, lastlogin, usertype FROM users WHERE id = $1", u.Id).
		Scan(
			&user.Id,
			&user.Completename,
			&user.Cpf,
			&user.Phone,
			&user.Email,
			&user.Password,
			&user.Country,
			&user.CreateAt,
			&user.LastLogin,
			&user.UserType,
		)
	if err != nil {
		log.Println("GetUserById: ", err)
		return user, errors.New("Not found user by id")
	}
	return user, nil
}
