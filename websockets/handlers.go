package websockets

import (
	"errors"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
)

func (s subscription) handleMessage(msg dto.WebSocketMessage) error {
	switch msg.Command {
	case models.CreateMessage:
		return s.handleCreateMessage(msg)
	case models.UpdateMessage:
		return s.handleUpdateMessage(msg)
	case models.DeleteMessage:
		return s.handleDeleteMessage(msg)
	case models.RemoveUser:
		return s.handleRemoveUser(msg)
	case models.AddUser:
		return s.handleAddUser(msg)
	case models.CreateRoom:
		return s.handleCreateRoom(msg)
	case models.UpdateRoom:
		return s.handleUpdateRoom(msg)
	case models.DeleteRoom:
		return s.handleDeleteRoom(msg)
	default:
		return errors.New("invalid message command")
	}
}

func (s subscription) handleCreateMessage(msg dto.WebSocketMessage) error {

	dto := dto.CreateMessage{
		Text: msg.Text,
	}

	if err := s.verifyUserInRoom(msg.RoomID); err != nil {
		return err
	}

	m, err := s.messageService.CreateNoValidation(msg.RoomID, s.userID, dto)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) verifyUserInRoom(roomID string) error {

	for _, room := range s.rooms {
		if room == roomID {
			return nil
		}
	}

	return errors.New("user not in room")
}

func (s subscription) handleUpdateMessage(msg dto.WebSocketMessage) error {

	dto := dto.UpdateMessage{
		Text: msg.Text,
	}

	m, err := s.messageService.Update(msg.TargetID, s.userID, dto)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleDeleteMessage(msg dto.WebSocketMessage) error {

	m, err := s.messageService.Delete(msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleRemoveUser(msg dto.WebSocketMessage) error {

	err := s.roomService.RemoveUser(msg.RoomID, msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	m, err := s.messageService.CreateRemoveUser(msg.RoomID, msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleAddUser(msg dto.WebSocketMessage) error {

	err := s.roomService.AddUser(msg.RoomID, msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	m, err := s.messageService.CreateAddUser(msg.RoomID, msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	return s.broadcastGlobally(m)
}

func (s subscription) handleCreateRoom(msg dto.WebSocketMessage) error {

	dto := dto.CreateRoom{
		Name: msg.Text,
	}

	room, err := s.roomService.Create(dto, s.userID)
	if err != nil {
		return err
	}

	err = s.addRoom(room.ID)
	if err != nil {
		return err
	}

	m, err := s.messageService.CreateCreateRoom(room.ID, s.userID)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleUpdateRoom(msg dto.WebSocketMessage) error {

	dto := dto.UpdateRoom{
		Name: msg.Text,
	}

	_, err := s.roomService.Update(msg.RoomID, s.userID, dto)
	if err != nil {
		return err
	}

	m, err := s.messageService.CreateUpdateRoom(msg.RoomID, s.userID)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleDeleteRoom(msg dto.WebSocketMessage) error {

	err := s.roomService.Delete(msg.RoomID, s.userID)
	if err != nil {
		return err
	}

	m := models.Message{
		Text:     "Room deleted",
		RoomID:   msg.RoomID,
		UserID:   s.userID,
		Command:  models.DeleteRoom,
		TargetID: msg.RoomID,
	}

	return s.broadcast(m)
}
