package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
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

	// jsonBody, _ := json.Marshal(reqBody)
	// reader := bytes.NewReader(jsonBody)

	// Why we need to create a new client for each connection:
	// because if a service is down, we don't want to block the whole application
	// so we create a new client for each connection
	tripService, err := grpc_clients.NewServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to close the client to avoid resource leaks!
	defer tripService.Close()

	// resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	// if err != nil {
	// 	log.Print(err)
	// }

	// defer resp.Body.Close()

	// var respBody any

	// if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
	// 	http.Error(w, "Failed to parse json data from trip service", http.StatusBadRequest)
	// 	return
	// }

	response := contracts.APIResponse{
		Data: "okay",
	}
	if err := writeJson(w, http.StatusCreated, response); err != nil {
		panic("error ")
	}
}
