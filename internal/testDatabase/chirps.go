package testDatabase

import "errors"

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(chirp string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		Body: chirp,
		ID:   newId,
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