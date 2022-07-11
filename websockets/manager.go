package websockets

import (
	"net/http"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/rs/zerolog/log"
)

type IManager interface {
	ServeWs(w http.ResponseWriter, r *http.Request, roomID string, userID string)
	DisconnectUserFromRoom(userID string, roomID string)
	DisconnectRoom(roomID string)
}

type Manager struct{}

var managerOnce sync.Once
var manager IManager

func GetManager() IManager {
	managerOnce.Do(func() {
		log.Info().Msg("Initializing websocket manager")
		initHub()
		manager = &Manager{}
	})
	return manager
}

// ServeWs handles websocket requests from the peer.
func (m *Manager) ServeWs(w http.ResponseWriter, r *http.Request, roomID string, userID string) {
	log.Info().Str("room_id", roomID).Str("user_id", userID).Msg("New websocket connection")
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error().Err(err).Str("room_id", roomID).Str("user_id", userID).Msg("Failed to upgrade websocket connection")
		return
	}

	c := &connection{send: make(chan models.Message, 256), ws: ws}
	s := subscription{c, roomID, userID}
	h.register <- s

	go s.writePump()
	go s.readPump()
}

func (m *Manager) DisconnectUserFromRoom(userID string, roomID string) {
	log.Debug().Str("room_id", roomID).Str("user_id", userID).Msg("Disconnecting user from room connection")

	connections := h.rooms[roomID]
	for c, id := range connections {
		if id == userID {
			h.unregister <- subscription{c, roomID, userID}
		}
	}
}

func (m *Manager) DisconnectRoom(roomID string) {
	log.Debug().Str("room_id", roomID).Msg("Disconnecting room connection")

	connections := h.rooms[roomID]
	for c, id := range connections {
		h.unregister <- subscription{c, roomID, id}
	}
}
