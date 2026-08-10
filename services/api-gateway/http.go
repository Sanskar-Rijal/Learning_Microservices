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

	// Why we need to create a new client for each connection:
	// because if a service is down, we don't want to block the whole application
	// so we create a new client for each connection
	tripService, err := grpc_clients.NewServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := r.Context()

	tripPreview, err := tripService.Client.PreviewTrip(ctx, reqBody.ToProto())
	if err != nil {
		log.Printf("Failed to preview a trip: %v", err)
		http.Error(w, "Failed to preview a trip", http.StatusInternalServerError)
		return
	}

	// Don't forget to close the client to avoid resource leaks!
	defer tripService.Close()

	response := contracts.APIResponse{
		Data: tripPreview,
	}
	if err := writeJson(w, http.StatusCreated, response); err != nil {
		panic("error ")
	}
}
