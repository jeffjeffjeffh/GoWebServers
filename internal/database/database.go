package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mutex *sync.RWMutex
}

type DBstructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
	RevokedTokens map[string]time.Time `json:"revokedTokens"`
}

func newDB(filename string) *DB {
	return &DB{
		path: filename,
		mutex: &sync.RWMutex{},
	}
}

func newDBstructure() DBstructure {
	return DBstructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
		RevokedTokens: map[string]time.Time{},
	}
}

func CreateDB(filename string) (*DB, error) {
	db := newDB(filename)
	dbStructure := newDBstructure()
	return db, db.writeDB(dbStructure)
}

func LoadDB(filename string) (*DB, error) {
	file, err := os.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) || len(file) == 0 {
		return CreateDB(filename)
	}

	db := newDB(filename)
	return db, err
}

// takes a dbStructure already loaded from CreateChirp
func (db *DB) writeDB(dbStructure DBstructure) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		log.Println(err)
		return err
	}

	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return DBstructure{}, err
	}
	
	dbStructure := DBstructure{}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		log.Println(err)
		return DBstructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) CheckTokenStatus(tokenStr string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, ok := dbStructure.RevokedTokens[tokenStr]
	if ok {
		return true, nil
	}

	return false, nil
}

func (db *DB) RevokeToken(tokenString string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.RevokedTokens[tokenString]
	if ok {
		err := errors.New("token already revoked")
		log.Println(err)
		return err
	}

	dbStructure.RevokedTokens[tokenString] = time.Now()
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}