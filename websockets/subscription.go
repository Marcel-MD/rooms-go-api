package websockets

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type subscription struct {
	userID         string
	roomID         string
	pubsub         *redis.PubSub
	ws             *websocket.Conn
	close          chan struct{}
	messages       chan models.Message
	messageService services.IMessageService
	roomService    services.IRoomService
}

func connect(userID, roomID string, ws *websocket.Conn) (*subscription, error) {
	var s *subscription

	pubsub := rdb.Subscribe(ctx, roomID)
	err := pubsub.Ping(ctx)
	if err != nil {
		return s, err
	}

	s = &subscription{
		userID:         userID,
		roomID:         roomID,
		pubsub:         pubsub,
		ws:             ws,
		close:          make(chan struct{}),
		messages:       make(chan models.Message),
		messageService: services.GetMessageService(),
		roomService:    services.GetRoomService(),
	}

	go func() {
		log.Info().Str("user_id", userID).Str("room_id", roomID).Msg("User listening to room")

		for {
			select {
			case msg, ok := <-pubsub.Channel():
				if !ok {
					log.Warn().Str("user_id", userID).Str("room_id", roomID).Msg("Pubsub channel closed")
					return
				}

				var m models.Message
				err := json.Unmarshal([]byte(msg.Payload), &m)
				if err != nil {
					log.Err(err).Str("user_id", userID).Str("room_id", roomID).Msg("Failed to unmarshal message")
					continue
				}

				s.messages <- m

			case <-s.close:
				log.Info().Str("user_id", userID).Str("room_id", roomID).Msg("User stopped listening to room")
				return
			}
		}
	}()

	return s, nil
}

func (s *subscription) disconnect() error {
	if s.pubsub == nil {
		return errors.New("subscriber is not connected")
	}

	if err := s.pubsub.Unsubscribe(ctx); err != nil {
		return err
	}

	if err := s.pubsub.Close(); err != nil {
		return err
	}

	s.close <- struct{}{}

	close(s.messages)

	return nil
}

func (s *subscription) broadcast(message models.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return rdb.Publish(ctx, s.roomID, payload).Err()
}

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
			log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Unexpected websocket close")
			break
		}

		if len(dto.Text) > 500 {
			log.Warn().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Invalid message text length")
			continue
		}

		if len(dto.Command) < 1 || len(dto.Command) > 50 {
			log.Warn().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Invalid message command length")
			continue
		}

		if len(dto.TargetID) < 1 || len(dto.TargetID) > 50 {
			log.Warn().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Invalid message target id length")
			continue
		}

		err = s.handleMessage(dto)
		if err != nil {
			log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Failed to handle message")
			break
		}
	}
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
		case message, ok := <-s.messages:
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
