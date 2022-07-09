package websockets

import (
	"log"
	"net/http"
	"time"

	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type connection struct {
	ws   *websocket.Conn
	send chan models.Message
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, roomID string, userID string) {
	log.Print(roomID)
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &connection{send: make(chan models.Message, 256), ws: ws}
	s := subscription{c, roomID, userID}
	h.register <- s

	go s.writePump()
	go s.readPump()
}
