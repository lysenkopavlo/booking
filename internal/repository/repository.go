package repository

import (
	"time"

	"github.com/lysenkopavlo/booking/internal/models"
)

type DataBaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(res models.RoomRestriction) error

	SearchAvailabilityByDatesAndRoomID(startDate, endDate time.Time, roomID int) (bool, error)

	SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error)

	GetRoomByID(roomID int) (models.Room, error)
}
