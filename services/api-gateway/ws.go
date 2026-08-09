package main

import (
	"fmt"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	// upgrade the connection http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID") // front will have to send ?userID=blah_blah_blah like this
	if userID == "" {
		log.Println("You must send userID on query")
		return
	}
	// using for loop listening for messages
	for {
		// read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("message received %s \n", message)

	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	// upgrade the connection http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("You must send userID on query")
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("You must send userID on query")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Println("You must send packageSlug on query")
		return
	}

	type Driver struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate       string `json:"carPlate"`
		PackageSlug    string `json:"packageSlug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			Id:             userID,
			Name:           "Sanskar",
			ProfilePicture: util.GetRandomAvatar(7),
			CarPlate:       "abcdefghi",
			PackageSlug:    packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}
