package models

import "time"

// User us a model for users_table
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccessLevel int
}

// Room is a model for rooms_table
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Reservation is a model for reservations table
type Reservation struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	StartDate   time.Time
	EndDate     time.Time
	RoomID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccessLevel int
	Room        Room
}

// Restriction is a model for restrictions table
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// RoomRestriction is a model for room_restriction table
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionId int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
