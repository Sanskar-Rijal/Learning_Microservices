package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse json data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// simple validation
	if reqBody.UserID == "" {
		http.Error(w, "Pass userID on body", http.StatusBadRequest)
		return
	}

	// Todo call trip service
	log.Println("SUCCESS")
	response := contracts.APIResponse{
		Data: "okeyyyy",
	}
	if err := writeJson(w, http.StatusCreated, response); err != nil {
		panic("error ")
	}
}
