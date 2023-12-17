package testDatabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mutex *sync.RWMutex
}

type DBstructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

type Chirp struct{
	ID int `json:"id"`
	Body string `json:"body"`
}

type User struct{
	Email string `json:"email"`
	ID int `json:"id"`
}

func NewDB(filename string) (*DB, error) {
	db := &DB{
		path: filename,
		mutex: &sync.RWMutex{},
	}

	err := db.ensureDB()

	return db, err
}

func (db *DB) ensureDB() error {
	file, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) || len(file) == 0 {
		fmt.Println("db file does not exist or is empty. creating new file...")
		return db.initialize()
	}
	if err == nil {
		fmt.Println("db file found.")
	}
	return err
}

func (db *DB) initialize() error {
	dbStructure := DBstructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) CreateChirp(chirp string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		Body: chirp,
		ID: newId,
	}

	dbStructure.Chirps[newId] = newChirp

	return newChirp, db.writeDB(dbStructure)
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}

	return chirp, nil
}

func (db *DB) ListChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := []Chirp{}
	for _, val := range dbStructure.Chirps {
		chirps = append(chirps, val)
	}

	return chirps, nil
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newId := len(dbStructure.Users) + 1
	user := User{
		ID: newId,
		Email: email,
	}

	dbStructure.Users[newId] = user

	return user, db.writeDB(dbStructure)
}

// takes a dbStructure already loaded from CreateChirp
func (db *DB) writeDB(dbStructure DBstructure) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

// used by CreateChirp and ReadChirps to load the file into a DBstructure
func (db *DB) loadDB() (DBstructure, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBstructure{}, err
	}
	
	dbStructure := DBstructure{}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBstructure{}, err
	}

	return dbStructure, nil
}