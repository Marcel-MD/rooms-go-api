package websockets

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Marcel-MD/rooms-go-api/logger"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type subscription struct {
	userID         string
	rooms          map[string]bool
	pubsub         *redis.PubSub
	ws             *websocket.Conn
	rdb            *redis.Client
	ctx            context.Context
	close          chan struct{}
	send           chan models.Message
	messageService services.IMessageService
	roomService    services.IRoomService
	userService    services.IUserService
}

func connect(userID string, rooms map[string]bool, ws *websocket.Conn, rdb *redis.Client, ctx context.Context) (*subscription, error) {

	roomIDs := make([]string, 0, len(rooms))
	for r := range rooms {
		roomIDs = append(roomIDs, r)
	}

	pubsub := rdb.Subscribe(ctx, roomIDs...)
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
		userService:    services.GetUserService(),
	}

	go s.listen()

	s.userService.SetIsOnline(userID, true)

	return s, nil
}

func (s *subscription) listen() {
	log.Info().Str(logger.UserID, s.userID).Msg("User listening to rooms")

	for {
		select {
		case msg, ok := <-s.pubsub.Channel():
			if !ok {
				log.Warn().Str(logger.UserID, s.userID).Msg("Pubsub channel closed")
				return
			}

			var m models.Message
			err := json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.Err(err).Str(logger.UserID, s.userID).Msg("Failed to unmarshal message")
				continue
			}

			s.send <- m

		case <-s.close:
			log.Info().Str(logger.UserID, s.userID).Msg("User stopped listening to rooms")
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

	s.userService.SetIsOnline(s.userID, false)

	return nil
}

func (s *subscription) reconnect() error {
	log.Info().Str(logger.UserID, s.userID).Msg("User reconnecting to rooms")

	if err := s.pubsub.Unsubscribe(s.ctx); err != nil {
		return err
	}

	if err := s.pubsub.Close(); err != nil {
		return err
	}

	s.close <- struct{}{}

	roomIDs := make([]string, 0, len(s.rooms))
	for r := range s.rooms {
		roomIDs = append(roomIDs, r)
	}

	pubsub := s.rdb.Subscribe(s.ctx, roomIDs...)
	err := pubsub.Ping(s.ctx)
	if err != nil {
		return err
	}

	s.pubsub = pubsub

	go s.listen()

	return nil
}

func (s *subscription) addRoom(room string) error {
	s.rooms[room] = true
	return s.reconnect()
}

func (s *subscription) removeRoom(room string) error {
	delete(s.rooms, room)
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
