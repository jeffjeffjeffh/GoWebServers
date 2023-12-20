package database

import (
	"errors"
	"log"
	"time"
)

func (db *DB) CheckTokenStatus(tokenStr string) (bool, error) {
	dbStructure, err := db.load()
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
	dbStructure, err := db.load()
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