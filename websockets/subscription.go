package websockets

import (
	"context"
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
	rooms          []string
	pubsub         *redis.PubSub
	ws             *websocket.Conn
	rdb            *redis.Client
	ctx            context.Context
	close          chan struct{}
	send           chan models.Message
	messageService services.IMessageService
	roomService    services.IRoomService
}

func connect(userID string, rooms []string, ws *websocket.Conn, rdb *redis.Client, ctx context.Context) (*subscription, error) {

	pubsub := rdb.Subscribe(ctx, rooms...)
	err := pubsub.Ping(ctx)
	if err != nil {
		return nil, err
	}

	s := &subscription{
		userID:         userID,
		rooms:          rooms,
		pubsub:         pubsub,
		ws:             ws,
		rdb:            rdb,
		ctx:            ctx,
		close:          make(chan struct{}),
		send:           make(chan models.Message),
		messageService: services.GetMessageService(),
		roomService:    services.GetRoomService(),
	}

	go s.listen()

	return s, nil
}

func (s *subscription) listen() {
	log.Info().Str("user_id", s.userID).Msg("User listening to rooms")

	for {
		select {
		case msg, ok := <-s.pubsub.Channel():
			if !ok {
				log.Warn().Str("user_id", s.userID).Msg("Pubsub channel closed")
				return
			}

			var m models.Message
			err := json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.Err(err).Str("user_id", s.userID).Msg("Failed to unmarshal message")
				continue
			}

			s.send <- m

		case <-s.close:
			log.Info().Str("user_id", s.userID).Msg("User stopped listening to rooms")
			return
		}
	}
}

func (s *subscription) disconnect() error {
	if s.pubsub == nil {
		return errors.New("subscriber is not connected")
	}

	if err := s.pubsub.Unsubscribe(s.ctx); err != nil {
		return err
	}

	if err := s.pubsub.Close(); err != nil {
		return err
	}

	s.close <- struct{}{}

	close(s.send)

	return nil
}

func (s *subscription) reconnect() error {
	log.Info().Str("user_id", s.userID).Msg("User reconnecting to rooms")

	if err := s.pubsub.Unsubscribe(s.ctx); err != nil {
		return err
	}

	if err := s.pubsub.Close(); err != nil {
		return err
	}

	s.close <- struct{}{}

	pubsub := s.rdb.Subscribe(s.ctx, s.rooms...)
	err := pubsub.Ping(s.ctx)
	if err != nil {
		return err
	}

	s.pubsub = pubsub

	go s.listen()

	return nil
}

func (s *subscription) addRoom(room string) error {
	s.rooms = append(s.rooms, room)
	return s.reconnect()
}

func (s *subscription) removeRoom(room string) error {
	for i, r := range s.rooms {
		if r == room {
			s.rooms = append(s.rooms[:i], s.rooms[i+1:]...)
			break
		}
	}

	return s.reconnect()
}

func (s *subscription) broadcast(message models.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.rdb.Publish(s.ctx, message.RoomID, payload).Err()
}

func (s *subscription) broadcastGlobally(message models.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.rdb.Publish(s.ctx, globalChannel, payload).Err()
}
