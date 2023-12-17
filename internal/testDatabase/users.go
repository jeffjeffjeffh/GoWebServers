package testDatabase

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Password []byte `json:"password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newId := len(dbStructure.Users) + 1

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:       newId,
		Email:    email,
		Password: hasedPassword,
	}

	dbStructure.Users[newId] = user

	return user, db.writeDB(dbStructure)
}