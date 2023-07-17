// Package dbrepo holds postgres queries

package dbrepo

import (
	"time"

	"github.com/lysenkopavlo/booking/internal/models"
)

func (m *testDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (m *testDbRepo) InsertReservation(res models.Reservation) (int, error) {
	return -1, nil
}

// InsertRoomRestriction inserts a room_restriction into the database
func (m *testDbRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDatesAndRoomID search available rooms by specific dates and room_id
// returns true if there is a free dates for specific room (aka room_id)
func (m *testDbRepo) SearchAvailabilityByDatesAndRoomID(startDate, endDate time.Time, roomID int) (bool, error) {
	// numRow holds the number of row
	// which crosses with startDate and endDate
	var numRow int

	return numRow == 0, nil
}

// SearchAvailabilityForAllRooms search available rooms by specific dates
// returns slice of models.Room, if there is available one
func (m *testDbRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	// Rooms holds the id and room_name for available rooms
	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID gets the room by its id
func (m *testDbRepo) GetRoomByID(roomID int) (models.Room, error) {
	var r models.Room

	return r, nil
}
