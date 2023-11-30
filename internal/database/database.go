package database

import (
	"encoding/json"
	// "fmt"
	"os"
	"sync"
	"syscall"
)

type dbConnection struct {
	path string
	mux *sync.RWMutex
}

type dbFormat struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func LoadFile(path string) (dbFormat, error) {
	conn := dbConnection{
		path: path,
		mux: &sync.RWMutex{},
	}

	file, err := os.ReadFile(conn.path)
	if err != nil {
		os.WriteFile("db.json", []byte("asdf"), syscall.S_IRUSR | syscall.S_IWUSR)
		file, _ = os.ReadFile(conn.path)
	}

	// fmt.Println(string(file))

	return loadDB(file)
}

func loadDB(file []byte) (dbFormat, error) {
	db := dbFormat{}
	err := json.Unmarshal(file, &db)
	if err != nil {
		return dbFormat{}, err
	}

	// fmt.Println(db)

	return db, nil
}

func createChirp() {
	
}