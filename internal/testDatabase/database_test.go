package testDatabase

import (
	"testing"
)

func TestWriteChirp(t *testing.T) {
	db := newDB("database.json")

	testChirp := "This is the first chirp ever!"

	_, err := db.CreateChirp(testChirp)
	if err != nil {		
		t.Error(err)
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		t.Error(err)
	}

	_, ok := dbStructure.Chirps[1]
	if !ok {
		t.Fail()
	}
}