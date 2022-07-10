package websockets

import (
	"sync"

	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/services"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	rooms      map[string]map[*connection]bool
	broadcast  chan models.Message
	register   chan subscription
	unregister chan subscription
	service    services.IMessageService
}

var once sync.Once
var h hub

func InitHub() {
	once.Do(func() {
		h = hub{
			rooms:      make(map[string]map[*connection]bool),
			broadcast:  make(chan models.Message),
			register:   make(chan subscription),
			unregister: make(chan subscription),
			service:    services.GetMessageService(),
		}
		go h.run()
	})
}

func (h *hub) run() {
	for {
		select {

		case s := <-h.register:
			connections := h.rooms[s.roomID]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.roomID] = connections
			}
			h.rooms[s.roomID][s.conn] = true

		case s := <-h.unregister:
			connections := h.rooms[s.roomID]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.roomID)
					}
				}
			}

		case m := <-h.broadcast:
			connections := h.rooms[m.RoomID]
			for c := range connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.RoomID)
					}
				}
			}
		}
	}
}
