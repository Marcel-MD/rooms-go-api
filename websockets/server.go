package websockets

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IServer interface {
	ServeWS(w http.ResponseWriter, r *http.Request, roomID, userID string) error
}

type wsServer struct{}

var (
	serverOnce sync.Once
	server     IServer
)

func GetServer() IServer {
	serverOnce.Do(func() {
		log.Info().Msg("Initializing websocket server")
		initRDB()
		server = &wsServer{}
	})
	return server
}

func (m *wsServer) ServeWS(w http.ResponseWriter, r *http.Request, roomID, userID string) error {
	log.Info().Str("room_id", roomID).Str("user_id", userID).Msg("New websocket connection")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Str("room_id", roomID).Str("user_id", userID).Msg("Failed to upgrade websocket connection")
		return err
	}

	s, err := connect(userID, roomID, ws)
	if err != nil {
		log.Error().Err(err).Str("room_id", roomID).Str("user_id", userID).Msg("Failed to connect to room")
		ws.Close()
		return err
	}

	go s.writePump()
	go s.readPump()

	return nil
}
