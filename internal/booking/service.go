package booking

import "context"

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{store}
}

func (s *Service) Book(b Booking) (Booking, error) {
	return s.store.Book(b)
}

func (s *Service) ListBookings(movieID string) []Booking {
	return s.store.ListBookings(movieID)
}

func (s *Service) ConfirmSeat(ctx context.Context, sessionID, userID string) (Booking, error) {
	return s.store.ConfirmSeat(ctx, sessionID, userID)
}

func (s *Service) ReleaseSeat(ctx context.Context, sessionID, userID string) error {
	return s.store.ReleaseSeat(ctx, sessionID, userID)
}
