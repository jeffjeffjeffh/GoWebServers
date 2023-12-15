package testDatabase

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mutex *sync.RWMutex
}

type DBstructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct{
	ID int `json:"id"`
	Body string `json:"body"`
}

func NewDB(filename string) (*DB, error) {
	db := DB{
		path: filename,
		mutex: &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return nil, nil
	}

	return &db, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.initialize()
	}
	return nil
}

func (db *DB) initialize() error {
	dbStructure := DBstructure{
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) CreateChirp(chirp string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	
	newId := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		Body: chirp,
		ID: newId,
	}

	dbStructure.Chirps[newId] = newChirp

	return db.writeDB(dbStructure)
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

	return dbStructure, nil
}