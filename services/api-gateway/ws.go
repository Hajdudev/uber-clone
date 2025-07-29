package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error making the websocker connection %v", err)
	}
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		println("No userID provided")
		return
	}
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading the message from websocker: %v", err)
			break
		}
		log.Printf("Received message: %v", message)
	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error making the websocker connection %v", err)
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		println("No userID provided")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		println("No packageSlug provided")
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
			Name:           "Hajdu",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "SSD234",
			PackageSlug:    packageSlug,
		},
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error writing to the websocket: %v", err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading the message from websocker: %v", err)
			break
		}
		log.Printf("Received message: %v", message)
	}
}
