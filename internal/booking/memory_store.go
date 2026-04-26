package booking

type MemoryStore struct {
	booking map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		booking: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) (Booking, error) {
	if _, exist := s.booking[b.SeatID]; exist {
		return Booking{}, ErrSeatAlreadyBooked
	}

	s.booking[b.SeatID] = b
	return b, nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var ResBookings []Booking
	for _, b := range s.booking {
		if movieID == b.MovieID {
			ResBookings = append(ResBookings, b)
		}
	}
	return ResBookings
}
