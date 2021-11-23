package service

import (
	"context"
	"fmt"
	"github.com/nure-mates/api/src/models"
)

func (s *Service) CreateRoom(ctx context.Context, room *models.Room) error {
	if err := s.roomRepo.CreateRoom(ctx, room); err != nil {
		return fmt.Errorf("create room %s by %d: %v", room.Name, room.HostID, err)
	}
	return nil
}

func (s *Service) GetRoom(ctx context.Context, id int) (*models.Room, error) {
	room, err := s.roomRepo.GetRoom(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get room %d: %v", id, err)
	}
	userCount := len(room.UserIDs)
	if err != nil {
		return nil, err
	}
	room.UserCount = userCount
	return room, nil
}

func (s *Service) GetUserRooms(ctx context.Context, userID int) ([]models.Room, error) {
	rooms, err := s.roomRepo.GetUserRooms(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user %d rooms: %v", userID, err)
	}
	for i := range rooms {
		c := len(rooms[i].UserIDs)
		if err != nil {
			return nil, err
		}
		rooms[i].UserCount = c
	}
	return rooms, nil
}

func (s *Service) AddUserToRoom(ctx context.Context, roomID, userID int) error {
	if err := s.roomRepo.AddUserToRoom(ctx, roomID, userID); err != nil {
		return fmt.Errorf("add user %d to room %d: %v", userID, roomID, err)
	}

	return nil
}

func (s *Service) RemoveUserFromRoom(ctx context.Context, roomID, userID int) error {
	if err := s.roomRepo.RemoveUserFromRoom(ctx, roomID, userID); err != nil {
		return fmt.Errorf("remove user %d from room %d: %v", userID, roomID, err)
	}
	if noUsers, err := s.CheckRoom(ctx, roomID); err != nil {
		if noUsers {
			err := s.DeleteRoom(ctx, roomID)
			return err
		}
	} else {
		return fmt.Errorf("remove empty room %d: %v", roomID, err)
	}
	return nil
}

func (s *Service) DeleteRoom(ctx context.Context, id int) error {
	if err := s.roomRepo.DeleteRoom(ctx, id); err != nil {
		return fmt.Errorf("delete room %d: %v", id, err)
	}

	return nil
}

func (s *Service) GetAvailableRooms(ctx context.Context, userID int) ([]models.Room, error) {
	rooms, err := s.roomRepo.GetAvailableRooms(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get available rooms for user %d: %v", userID, err)
	}

	return rooms, nil
}

func (s *Service) UpdateRoom(ctx context.Context, room *models.Room) error {
	if err := s.roomRepo.UpdateRoom(ctx, room); err != nil {
		return fmt.Errorf("update room %s, %d: %v", room.Name, room.ID, err)
	}

	return nil
}

func (s *Service) CheckRoom(ctx context.Context, roomID int) (bool, error) {
	empty, err := s.roomRepo.CheckRoom(ctx, roomID)
	return empty, err
}
