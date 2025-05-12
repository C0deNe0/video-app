package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllRoom RoomMap

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "create room request Handler")
	w.Header().Set("Access-control-Allow-Origin", "*")
	roomID := AllRoom.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}
	json.NewEncoder(w).Encode(resp{RoomID: roomID})

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast

		for _, client := range AllRoom.Map[msg.RoomID] {
			if client.conn != msg.Client {
				err := client.conn.WriteJSON(msg.Message)
				if err != nil {
					log.Fatal(err)
					client.conn.Close()
				}
			}
		}
	}
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
	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("Read Error: ", err)

		}

		msg.Client = ws
		msg.RoomID = roomID[0]
		broadcast <- msg
	}
}
