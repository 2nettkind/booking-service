package booking

import (
	"sync"
)

type ConcurrentStore struct {
	booking map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		booking: map[string]Booking{},
	}
}

func (s *ConcurrentStore) Book(b Booking) (Booking, error) {
	s.Lock()
	defer s.Unlock()

	if _, exist := s.booking[b.SeatID]; exist {
		return Booking{}, ErrSeatAlreadyBooked
	}

	s.booking[b.SeatID] = b
	return b, nil
}

func (s *ConcurrentStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()

	var ResBookings []Booking
	for _, b := range s.booking {
		if movieID == b.MovieID {
			ResBookings = append(ResBookings, b)
		}
	}
	return ResBookings
}
