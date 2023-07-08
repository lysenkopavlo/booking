package repository

import "github.com/lysenkopavlo/booking/internal/models"

type DataBaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(res models.RoomRestriction) error
}