package server

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// single entity in the Hashmap
type Participant struct {
	Host bool
	conn *websocket.Conn
}

// RoomMap is the main hashmap which maps
// [roomID] ==> [[]participant]
type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}

// Init will initialize the RoomMap struct
func (r *RoomMap) Init() {
	//here we didn't need to setup the Mutex, bcoz the Hashmap will be created before any connection is started
	r.Map = make(map[string][]Participant)
}

// GetRoom will return the array of participants in the room
func (r *RoomMap) GetRoom(roomId string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.Unlock()

	return r.Map[roomId]

}

// create a uniqueID and return it => insert it in the hashmap
func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	roomID := string(b)
	r.Map[roomID] = []Participant{}

	return roomID

}

// InsertIntoRoom will create a participant add into the room
func (r *RoomMap) InsertIntoRoom(roomID string, host bool, coon *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{Host: host, conn: coon}

	r.Map[roomID] = append(r.Map[roomID], p)

}

// DeleteRoom is used to delete the Room
func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)
}
