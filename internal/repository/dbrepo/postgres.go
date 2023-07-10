package dbrepo

import (
	"context"
	"time"

	"github.com/lysenkopavlo/booking/internal/models"
)

func (m *postgresDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Making a variable to hold last id value
	// that returned from query
	var lastId int

	// stmt is a sql query with returning function
	stmt := `INSERT INTO reservations 
		(first_name, last_name, email, phone, start_date, end_date, 
		room_id, created_at, updated_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&lastId)

	if err != nil {
		return 0, err
	}
	return lastId, nil
}

// InsertRoomRestriction inserts a room_restriction into the database
func (m *postgresDbRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// stmt is a sql query
	stmt := `INSERT INTO room_restrictions 
		(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.RestrictionId,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDatesAndRoomID search available rooms by specific dates and room_id
// returns true if there is a free dates for specific room (aka room_id)
func (m *postgresDbRepo) SearchAvailabilityByDatesAndRoomID(startDate, endDate time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// numRow holds the number of row
	//which crosses with startDate and endDate
	var numRow int

	//query holds SQL-query
	query := `SELECT FROM room_restriction
			COUNT(id)
		WHERE
			room_id = $1
			$2 < end_date AND $3 > start_date;`

	// Executing SQLquery with specified dates and room_id
	row := m.DB.QueryRowContext(ctx, query, startDate, endDate, roomID)
	err := row.Scan(&numRow)
	if err != nil {
		return false, err
	}
	return numRow == 0, nil
}
