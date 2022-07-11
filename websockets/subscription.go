package websockets

import (
	"time"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type subscription struct {
	conn   *connection
	roomID string
	userID string
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Starting websocket read pump")

	c := s.conn
	defer func() {
		log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Stopping websocket read pump")
		h.unregister <- s
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var dto dto.CreateMessage
		err := c.ws.ReadJSON(&dto)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Normal websocket close")
				break
			}
			log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Unexpected websocket close")
			break
		}

		if len(dto.Text) < 1 || len(dto.Text) > 500 {
			log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Invalid message length")
			continue
		}

		m, err := h.service.Create(dto, s.roomID, s.userID)
		if err != nil {
			log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Failed to create message")
			break
		}

		h.broadcast <- m
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Starting websocket write pump")

	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Stopping websocket write pump")
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Info().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Closing websocket connection")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.writeJSON(message); err != nil {
				log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Failed to write message")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Failed to write ping")
				return
			}
		}
	}
}
