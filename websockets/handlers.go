package websockets

import (
	"errors"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
)

func (s subscription) handleMessage(msg dto.WebSocketMessage) error {
	switch msg.Command {
	case models.Create:
		return s.handleCreate(msg)
	case models.Update:
		return s.handleUpdate(msg)
	case models.Delete:
		return s.handleDelete(msg)
	case models.RemoveUser:
		return s.handleRemoveUser(msg)
	case models.AddUser:
		return s.handleAddUser(msg)
	default:
		return errors.New("invalid message command")
	}
}

func (s subscription) handleCreate(msg dto.WebSocketMessage) error {

	dto := dto.CreateMessage{
		Text: msg.Text,
	}

	m, err := s.messageService.Create(s.roomID, s.userID, dto)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleUpdate(msg dto.WebSocketMessage) error {

	dto := dto.UpdateMessage{
		Text: msg.Text,
	}

	m, err := s.messageService.Update(msg.TargetID, s.userID, dto)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleDelete(msg dto.WebSocketMessage) error {

	m, err := s.messageService.Delete(msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleRemoveUser(msg dto.WebSocketMessage) error {

	err := s.roomService.RemoveUser(s.roomID, msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	m, err := s.messageService.CreateRemoveUser(s.roomID, msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}

func (s subscription) handleAddUser(msg dto.WebSocketMessage) error {

	err := s.roomService.AddUser(s.roomID, msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	m, err := s.messageService.CreateAddUser(s.roomID, msg.TargetID, s.userID)
	if err != nil {
		return err
	}

	return s.broadcast(m)
}
