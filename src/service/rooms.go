package service

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nure-mates/api/src/models"
)

func (s *Service) CreateRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	var err error

	if !*room.Public {
		room.Token, err = generateRandomString(models.TokenLenForRoom)
		if err != nil {
			log.Errorf("failed to create token: %v", err)

			return &models.Room{}, err
		}
	}

	err = s.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		log.Errorf("create room %s by %d: %v", room.Name, room.HostID, err)

		return &models.Room{}, err
	}

	return room, nil
}

func (s *Service) GetRoom(ctx context.Context, id int, userID int) (*models.Room, error) {
	room, err := s.roomRepo.GetRoom(ctx, id)
	if err != nil {
		log.Errorf("get room %d: %v", id, err)

		return nil, err
	}

	if room.HostID != userID {
		room.Token = ""
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

		// check if host
		if rooms[i].HostID != userID {
			rooms[i].Token = ""
		}
	}

	return rooms, nil
}

func (s *Service) AddUserToRoom(ctx context.Context, roomID, userID int) error {
	room, err := s.roomRepo.GetRoom(ctx, roomID)
	if err != nil {
		log.Errorf("add user %d to room %d: %v", userID, roomID, err)
		return err
	}

	if room.MaxUsers <= room.UserCount+1 {
		log.Errorf("max user count in %d room", room.ID)
		return err
	}

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
			err = s.DeleteRoom(ctx, roomID)

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
	if room.Public != nil && !*room.Public {
		token, err := generateRandomString(models.TokenLenForRoom)
		if err != nil {
			log.Errorf("failed to create token: %v", err)

			return err
		}

		room.Token = token
	}

	if room.Public != nil && *room.Public {
		room.Token = ""
	}

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

func (s *Service) JoinToRoom(ctx context.Context, token string, userID int) (*models.Room, error) {
	roomID, err := s.roomRepo.GetRoomIDViaToken(ctx, token)
	if err != nil {
		return &models.Room{}, err
	}

	userExists := false

	usersRoom, err := s.roomRepo.GetUsersInRoom(ctx, roomID)
	if err != nil {
		return &models.Room{}, err
	}

	for _, userRoom := range usersRoom {
		if userRoom.UserID == userID {
			userExists = true
			break
		}
	}

	if !userExists {
		err = s.AddUserToRoom(ctx, roomID, userID)
		if err != nil {
			return &models.Room{}, err
		}
	}

	resp, err := s.GetRoom(ctx, roomID, userID)
	if err != nil {
		return &models.Room{}, err
	}

	return resp, nil
}
