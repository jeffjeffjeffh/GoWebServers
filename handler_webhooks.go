package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	hook, err := decodeWebhooksParams(r)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if *hook.Event != "user.upgraded" {
		log.Println("ignoring webhook")
		w.WriteHeader(http.StatusOK)
		return
	}

	err = cfg.db.UpdateUserMembership(*hook.Data.UserID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	log.Println("user's chirpy red status updated")
	w.WriteHeader(http.StatusOK)
}