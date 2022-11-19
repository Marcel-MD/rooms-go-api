package websockets

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Marcel-MD/rooms-go-api/logger"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/rdb"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/go-redis/redis/v9"
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
	rdb         *redis.Client
	ctx         context.Context
}

var (
	serverOnce sync.Once
	server     IServer
)

func GetServer() IServer {
	serverOnce.Do(func() {
		log.Info().Msg("Initializing websocket server")
		rdb, ctx := rdb.GetRDB()

		server = &wsServer{
			userService: services.GetUserService(),
			rdb:         rdb,
			ctx:         ctx,
		}
	})
	return server
}

func (wss *wsServer) ServeWS(w http.ResponseWriter, r *http.Request, userID string) error {
	log.Info().Str(logger.UserID, userID).Msg("New websocket connection")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Str(logger.UserID, userID).Msg("Failed to upgrade websocket connection")
		return err
	}

	user, err := wss.userService.FindOne(userID)
	if err != nil {
		return err
	}

	rooms := user.Rooms
	roomIDs := make(map[string]bool)
	for _, r := range rooms {
		roomIDs[r.ID] = true
	}

	roomIDs[globalChannel] = true
	roomIDs[models.GeneralRoomID] = true
	roomIDs[models.AnnouncementsRoomID] = true

	s, err := connect(userID, roomIDs, ws, wss.rdb, wss.ctx)
	if err != nil {
		log.Error().Err(err).Str(logger.UserID, userID).Msg("Failed to connect to room")
		ws.Close()
		return err
	}

	go s.writePump()
	go s.readPump()

	return nil
}
