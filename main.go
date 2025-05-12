package main

import (
	"log"
	"net/http"

	"github.com/C0deNe0/video-chat-app/server"
)

func main() {

	server.AllRoom.Init()
	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	log.Println("starting server on port 8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
