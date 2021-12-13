package postgres

import (
	"context"

	"github.com/nure-mates/api/src/models"
)

type RoomRepo struct {
	*Postgres
}

func (p *Postgres) NewRoomRepo() *RoomRepo {
	p.DB.RegisterModel((*models.UsersRooms)(nil))
	return &RoomRepo{p}
}

func (r *RoomRepo) CreateRoom(ctx context.Context, room *models.Room) error {
	_, err := r.DB.NewInsert().
		Model(room).
		Returning("id").
		Exec(ctx)

	if err != nil {
		return err
	}

	err = r.AddUserToRoom(ctx, room.ID, room.HostID)

	return err
}

func (r *RoomRepo) GetUserRooms(ctx context.Context, userID int) ([]models.Room, error) {
	var rooms []models.Room

	if err := r.DB.NewSelect().
		Model(&rooms).
		Where("host_id = ?", userID).
		OrderExpr("id ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *RoomRepo) GetRoom(ctx context.Context, id int) (*models.Room, error) {
	var room models.Room
	if err := r.DB.NewSelect().
		Model(&room).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *RoomRepo) AddUserToRoom(ctx context.Context, roomID, userID int) error {
	rel := models.UsersRooms{RoomID: roomID, UserID: userID}

	_, err := r.DB.NewInsert().
		Model(&rel).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoomRepo) RemoveUserFromRoom(ctx context.Context, roomID, userID int) error {
	rel := models.UsersRooms{
		UserID: userID,
		RoomID: roomID,
	}

	_, err := r.DB.NewDelete().
		Model(&rel).
		Where("user_id = ? AND room_id = ?", userID, roomID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoomRepo) DeleteRoom(ctx context.Context, id int) error {
	_, err := r.DB.NewDelete().
		Model((*models.Room)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoomRepo) GetAvailableRooms(ctx context.Context, userID int) ([]models.Room, error) {
	var (
		rooms      []models.Room
		usersRooms []models.UsersRooms
	)

	if err := r.DB.NewSelect().
		Model(&usersRooms).
		Relation("Room").
		Where("user_id = ?", userID).
		Scan(ctx); err != nil {
		return nil, err
	}

	for i := range usersRooms {
		rooms = append(rooms, *usersRooms[i].Room)
	}

	return rooms, nil
}

func (r *RoomRepo) UpdateRoom(ctx context.Context, room *models.Room) error {
	if _, err := r.DB.NewUpdate().
		Model(room).
		OmitZero().
		Where("id = ?", room.ID).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *RoomRepo) CheckRoom(ctx context.Context, id int) (bool, error) {
	userCount, err := r.DB.NewSelect().Model((*models.UsersRooms)(nil)).Where("id = ?", id).Count(ctx)
	if err != nil {
		return false, err
	}

	if userCount > 0 {
		return false, nil
	}

	return true, nil
}

func (r *RoomRepo) GetUsersInRoom(ctx context.Context, id int) ([]models.UsersRooms, error) {
	var users []models.UsersRooms

	err := r.DB.NewSelect().Model(&users).Relation("User").Where("room_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *RoomRepo) GetPublicRooms(ctx context.Context) ([]models.Room, error) {
	var rooms []models.Room

	err := r.DB.NewSelect().Model(&rooms).Where("public = ?", true).Scan(ctx)

	return rooms, err
}

func (r *RoomRepo) GetUser(ctx context.Context, id int) (*models.User, error) {
	var user *models.User

	err := r.DB.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
