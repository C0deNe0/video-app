package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllRoom RoomMap

type resp struct {
	roomID string `json:"room_id"`
}

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "create room request Handler")
	roomID := AllRoom.CreateRoom()

	json.NewEncoder(w).Encode(resp{roomID: roomID})

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "join room request Handler")

	roomID, ok := r.URL.Query()["roomID"]
	if !ok {
		fmt.Println("RoomID missing in URL")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Websocket upgrader Error", err)
	}

	AllRoom.InsertIntoRoom(roomID[0], false, ws)

}
