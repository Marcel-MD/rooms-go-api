package websockets

import (
	"encoding/json"
	"errors"

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
	send           chan models.Message
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
		send:           make(chan models.Message),
		messageService: services.GetMessageService(),
		roomService:    services.GetRoomService(),
	}

	go s.listen()

	return s, nil
}

func (s *subscription) listen() {
	log.Info().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("User listening to room")

	for {
		select {
		case msg, ok := <-s.pubsub.Channel():
			if !ok {
				log.Warn().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Pubsub channel closed")
				return
			}

			var m models.Message
			err := json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.Err(err).Str("user_id", s.userID).Str("room_id", s.roomID).Msg("Failed to unmarshal message")
				continue
			}

			s.send <- m

		case <-s.close:
			log.Info().Str("user_id", s.userID).Str("room_id", s.roomID).Msg("User stopped listening to room")
			return
		}
	}
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

	close(s.send)

	return nil
}

func (s *subscription) broadcast(message models.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return rdb.Publish(ctx, s.roomID, payload).Err()
}
