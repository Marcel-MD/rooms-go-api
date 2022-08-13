package websockets

import (
	"net/http"
	"sync"
	"time"

	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
	globalChannel  = "global"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IServer interface {
	ServeWS(w http.ResponseWriter, r *http.Request, userID string) error
}

type wsServer struct {
	userService services.IUserService
}

var (
	serverOnce sync.Once
	server     IServer
)

func GetServer() IServer {
	serverOnce.Do(func() {
		log.Info().Msg("Initializing websocket server")
		initRDB()
		server = &wsServer{
			userService: services.GetUserService(),
		}
	})
	return server
}

func (wss *wsServer) ServeWS(w http.ResponseWriter, r *http.Request, userID string) error {
	log.Info().Str("user_id", userID).Msg("New websocket connection")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to upgrade websocket connection")
		return err
	}

	user, err := wss.userService.FindOne(userID)
	if err != nil {
		return err
	}

	rooms := user.Rooms
	roomsID := make([]string, len(rooms)+1)
	for _, r := range rooms {
		roomsID = append(roomsID, r.ID)
	}

	roomsID = append(roomsID, globalChannel)

	s, err := connect(userID, roomsID, ws)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to connect to room")
		ws.Close()
		return err
	}

	go s.writePump()
	go s.readPump()

	return nil
}
