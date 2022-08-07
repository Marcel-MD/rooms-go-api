package websockets

import (
	"errors"
	"time"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func (s subscription) readPump() {
	log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Starting websocket read pump")

	defer func() {
		log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Stopping websocket read pump")
		s.disconnect()
		s.ws.Close()
	}()

	s.ws.SetReadLimit(maxMessageSize)
	s.ws.SetReadDeadline(time.Now().Add(pongWait))
	s.ws.SetPongHandler(func(string) error { s.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var dto dto.WebSocketMessage
		err := s.ws.ReadJSON(&dto)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Normal websocket close")
				break
			}
			s.writeError(err)
			log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Unexpected websocket close")
			break
		}

		err = validateMessage(dto)
		if err != nil {
			s.writeError(err)
			log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Invalid message")
			continue
		}

		err = s.handleMessage(dto)
		if err != nil {
			s.writeError(err)
			log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Failed to handle message")
			continue
		}
	}
}

func validateMessage(message dto.WebSocketMessage) error {
	if len(message.Text) < 1 || len(message.Text) > 500 {
		return errors.New("invalid message text length")
	}

	if len(message.Command) < 1 || len(message.Command) > 50 {
		return errors.New("invalid message command length")
	}

	if len(message.TargetID) < 1 || len(message.TargetID) > 50 {
		return errors.New("invalid message target id length")
	}

	return nil
}

func (s *subscription) writePump() {
	log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Starting websocket write pump")

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.Debug().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Stopping websocket write pump")
		ticker.Stop()
		s.disconnect()
		s.ws.Close()
	}()

	for {
		select {
		case message, ok := <-s.send:
			if !ok {
				log.Info().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Closing websocket connection")
				s.write(websocket.CloseMessage, []byte{})
				return
			}

			if message.Command == models.RemoveUser && message.TargetID == s.userID {
				log.Info().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("User left the room")
				s.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := s.writeJSON(message); err != nil {
				log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Failed to write message")
				s.write(websocket.CloseMessage, []byte{})
				return
			}

		case <-ticker.C:
			if err := s.write(websocket.PingMessage, []byte{}); err != nil {
				log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Failed to write ping")
				return
			}
		}
	}
}

func (s *subscription) write(mt int, payload []byte) error {
	s.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return s.ws.WriteMessage(mt, payload)
}

func (s *subscription) writeJSON(v interface{}) error {
	s.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return s.ws.WriteJSON(v)
}

func (s *subscription) writeError(err error) error {
	s.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return s.ws.WriteJSON(map[string]string{"command": models.Error, "error": err.Error()})
}
