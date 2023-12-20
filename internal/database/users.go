package database

import (
	"errors"
	"log"

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

	_, ok := db.findUserByEmail(email)
	if ok {
		return User{}, errors.New("a user with that email already exists")
	}

	newId := len(dbStructure.Users) + 1

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:       newId,
		Email:    email,
		Password: hashedPassword,
	}

	dbStructure.Users[newId] = user

	return user, db.writeDB(dbStructure)
}

func (db *DB) UpdateUser(email, password string, id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return User{}, err
	}

	newUser := User{
		Email: email,
		Password: []byte(hashedPassword),
		ID: id,
	}

	_, ok := dbStructure.Users[id]
	if !ok {
		log.Println("user id not found in database")
		return User{}, errors.New("id not found")
	}

	dbStructure.Users[id] = newUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) Login(email, password string) (User, error) {
	log.Println("attempting login")

	user, ok := db.findUserByEmail(email)
	if !ok {
		return User{}, errors.New("user not found")
	}
	log.Println("user found")

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return User{}, errors.New("wrong password")
	}
	log.Println("password verified")


	return user, nil
}

func (db *DB) findUserByEmail(email string) (User, bool) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, false
	}

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, true
		}
	}

	return User{}, false
}