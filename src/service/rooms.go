package service

import (
	"context"
	"fmt"
	"github.com/nure-mates/api/src/models"
	log "github.com/sirupsen/logrus"
)

func (s *Service) CreateRoom(ctx context.Context, room *models.Room) error {
	if err := s.roomRepo.CreateRoom(ctx, room); err != nil {
		log.Errorf("create room %s by %d: %v", room.Name, room.HostID, err)
		return err
	}
	return nil
}

func (s *Service) GetRoom(ctx context.Context, id int) (*models.Room, error) {
	room, err := s.roomRepo.GetRoom(ctx, id)
	if err != nil {
		log.Errorf("get room %d: %v", id, err)
		return nil, err
	}
	usersRooms, err := s.getUsersInRoom(ctx, id)
	if err != nil {
		log.Errorf("get users in room %d: %v", id, err)
		return nil, err
	}
	for i := range usersRooms {
		user := usersRooms[i].User
		room.Users = append(room.Users, *user)
	}

	room.UserCount = len(room.Users)
	return room, nil
}

func (s *Service) GetUserRooms(ctx context.Context, userID int) ([]models.Room, error) {
	rooms, err := s.roomRepo.GetUserRooms(ctx, userID)
	if err != nil {
		log.Errorf("get user %d rooms: %v", userID, err)
		return nil, err
	}
	for i := range rooms {

		usersRoom, err := s.getUsersInRoom(ctx, rooms[i].ID)
		if err != nil {
			log.Errorf("get users in room %d: %v", rooms[i].ID, err)
			return nil, err
		}
		for j := range usersRoom {
			user := usersRoom[j].User
			rooms[i].Users = append(rooms[i].Users, *user)
		}
		rooms[i].UserCount = len(rooms[i].Users)

	}
	return rooms, nil
}

func (s *Service) AddUserToRoom(ctx context.Context, roomID, userID int) error {
	if err := s.roomRepo.AddUserToRoom(ctx, roomID, userID); err != nil {
		log.Errorf("add user %d to room %d: %v", userID, roomID, err)
		return err
	}

	return nil
}

func (s *Service) RemoveUserFromRoom(ctx context.Context, roomID, userID int) error {
	if err := s.roomRepo.RemoveUserFromRoom(ctx, roomID, userID); err != nil {
		log.Errorf("remove user %d from room %d: %v", userID, roomID, err)
		return err
	}
	if noUsers, err := s.CheckRoom(ctx, roomID); err != nil {
		log.Errorf("check room %d: %v", roomID, err)
		return err
	} else {
		if noUsers {
			err := s.DeleteRoom(ctx, roomID)
			return err
		}
	}
	return nil
}

func (s *Service) DeleteRoom(ctx context.Context, id int) error {
	if err := s.roomRepo.DeleteRoom(ctx, id); err != nil {
		log.Errorf("delete room %d: %v", id, err)
		return fmt.Errorf("delete room %d: %v", id, err)
	}

	return nil
}

func (s *Service) GetAvailableRooms(ctx context.Context, userID int) ([]models.Room, error) {
	rooms, err := s.roomRepo.GetAvailableRooms(ctx, userID)
	if err != nil {
		log.Errorf("get available rooms for user %d: %v", userID, err)
		return nil, err
	}

	return rooms, nil
}

func (s *Service) UpdateRoom(ctx context.Context, room *models.Room) error {
	if err := s.roomRepo.UpdateRoom(ctx, room); err != nil {
		log.Errorf("update room %s, %d: %v", room.Name, room.ID, err)
		return err
	}

	return nil
}

func (s *Service) CheckRoom(ctx context.Context, roomID int) (bool, error) {
	empty, err := s.roomRepo.CheckRoom(ctx, roomID)
	return empty, err
}

func (s *Service) GetPublicRooms(ctx context.Context) ([]models.Room, error) {
	rooms, err := s.roomRepo.GetPublicRooms(ctx)
	return rooms, err
}

func (s *Service) getUsersInRoom(ctx context.Context, roomID int) ([]models.UsersRooms, error) {
	users, err := s.roomRepo.GetUsersInRoom(ctx, roomID)
	return users, err
}
