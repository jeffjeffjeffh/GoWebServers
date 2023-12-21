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
	ChirpyRed bool	`json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.load()
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
		ChirpyRed: false,
	}

	dbStructure.Users[newId] = user

	return user, db.writeDB(dbStructure)
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

func (db *DB) UpdateUser(email, password string, id int) (User, error) {
	dbStruct, err := db.load()
	if err != nil {
		return User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return User{}, err
	}

	oldUser, ok := dbStruct.Users[id]
	if !ok {
		err := errors.New("id not found")
		log.Println(err)
		return User{}, err 
	}

	updatedUser := User{
		ID: oldUser.ID,
		Email: email,
		Password: hashedPassword,
		ChirpyRed: oldUser.ChirpyRed,
	}

	dbStruct.Users[id] = updatedUser

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (db *DB) UpdateUserMembership(id int) error {
	dbStruct, err := db.load()
	if err != nil {
		log.Println(err)
		return err
	}

	oldUser, ok := dbStruct.Users[id]
	if !ok {
		err := errors.New("user not found")
		log.Println(err)
		return err
	}

	updatedUser := User{
		Email: oldUser.Email,
		ID: oldUser.ID,
		Password: oldUser.Password,
		ChirpyRed: true,
	}

	dbStruct.Users[id] = updatedUser

	err = db.writeDB(dbStruct)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}	

func (db *DB) findUserByEmail(email string) (User, bool) {
	dbStruct, err := db.load()
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