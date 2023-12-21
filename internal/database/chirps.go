package database

import (
	"errors"
	"log"
)

type Chirp struct{
	AuthorID int `json:"author_id"`
	Body string `json:"body"`
	ID   int    `json:"id"`
}

func (db *DB) CreateChirp(chirp string, id int) (Chirp, error) {
	dbStructure, err := db.load()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		Body: chirp,
		ID:   newId,
		AuthorID: id,
	}

	dbStructure.Chirps[newId] = newChirp

	log.Println("chirp created")
	return newChirp, db.writeDB(dbStructure)
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.load()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}

	return chirp, nil
}

func (db *DB) ListChirps(authorId *int) ([]Chirp, error) {
	dbStruct, err := db.load()
	if err != nil {
		return []Chirp{}, err
	}
	
	chirps := []Chirp{}

	if authorId != nil {
		for _, val := range dbStruct.Chirps {
			if *authorId == val.AuthorID {
				chirps = append(chirps, val)
			}
		}
	} else {
		for _, val := range dbStruct.Chirps {
			chirps = append(chirps, val)
		}
	}

	return chirps, nil
}

func (db *DB) DeleteChirp(authorId, chirpId int) error {
	dbStructure, err := db.load()
	if err != nil {
		return err
	}

	chirpToDelete, ok := dbStructure.Chirps[chirpId]
	if !ok {
		err := errors.New("chirp not found")
		log.Println(err)
		return err
	}

	if chirpToDelete.AuthorID != authorId {
		err := errors.New("author id does not match")
		log.Println(err)
		return err
	}

	delete(dbStructure.Chirps, chirpId)
	return nil
}