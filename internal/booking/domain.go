package booking

import (
	"context"
	"errors"
	"time"
)

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	ConfirmSeat(ctx context.Context, sessionID, userID string) (Booking, error)
	ReleaseSeat(ctx context.Context, sessionID, userID string) error
}

var (
	ErrSeatAlreadyBooked = errors.New("this seat already taken")
)
