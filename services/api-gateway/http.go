package main

import (
	"bytes"
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

	jsonBody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

	resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()

	var respBody any

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, "Failed to parse json data from trip service", http.StatusBadRequest)
		return
	}

	response := contracts.APIResponse{
		Data: respBody,
	}
	if err := writeJson(w, http.StatusCreated, response); err != nil {
		panic("error ")
	}
}
