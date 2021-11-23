package postgres

import (
	"context"
	"github.com/nure-mates/api/src/models"
)

type RoomRepo struct {
	*Postgres
}

func (p *Postgres) NewRoomRepo() *RoomRepo {
	return &RoomRepo{p}
}

func (r *RoomRepo) CreateRoom(ctx context.Context, room *models.Room) error {
	_, err := r.DB.NewInsert().
		Model(room).
		Exec(ctx)

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
	rel := models.UserToRoom{RoomID: roomID, UserID: userID}
	_, err := r.DB.NewInsert().
		Model(rel).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepo) RemoveUserFromRoom(ctx context.Context, roomID, userID int) error {
	rel := models.UserToRoom{
		UserID: userID,
		RoomID: roomID,
	}
	_, err := r.DB.NewDelete().
		Model(rel).
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
		Where("room_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepo) GetAvailableRooms(ctx context.Context, userID int) ([]models.Room, error) {
	var rooms []models.Room

	if err := r.DB.NewSelect().
		Model((*models.Room)(nil)).
		Relation("").
		Where("user_id = ?", userID).
		Scan(ctx, &rooms); err != nil {
			return nil, err
	}
	return rooms, nil
}

func (r *RoomRepo) UpdateRoom(ctx context.Context, room *models.Room) error {
	if _, err := r.DB.NewUpdate().
		Model(&room).
		Where("id = ?", room.ID).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *RoomRepo) CheckRoom(ctx context.Context, id int) (bool, error) {
	userCount, err := r.DB.NewSelect().Model(&models.UserToRoom{}).ScanAndCount(ctx)
	if err != nil {
		return false, err
	}

	if userCount > 0 {
		return false, nil
	}
	return true, nil
}

