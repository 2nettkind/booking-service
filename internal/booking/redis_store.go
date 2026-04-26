package booking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute

// RedisStore implements session-based seat booking backed by Redis.
//
// Key design:
//
//	seat:{movieID}:{seatID}   → session JSON (TTL = held, no TTL = confirmed)
//	session:{sessionID}       → seat key     (reverse lookup)
type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (s *RedisStore) Book(b Booking) (Booking, error) {
	session, err := s.hold(b)
	if err != nil {
		return Booking{}, err
	}

	log.Printf("Session booked %v", session)

	return session, nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	ctx := context.Background()

	pattern := fmt.Sprintf("seat:%s:*", movieID)
	keys, err := s.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil
	}

	var bookings []Booking
	for _, key := range keys {
		val, err := s.rdb.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}
		var b Booking
		if err := json.Unmarshal(val, &b); err == nil {
			bookings = append(bookings, b)
		}
	}
	return bookings
}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)

	b.ID = id
	val, _ := json.Marshal(b)

	res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})
	ok := res.Val() == "OK"
	if !ok {
		return Booking{}, ErrSeatAlreadyBooked
	}

	s.rdb.Set(ctx, sessionKey(id), key, defaultHoldTTL)

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}

func (s *RedisStore) ReleaseSeat(ctx context.Context, sessionID, userID string) error {
	// Шаг 1: находим ключ места через сессию
	seatKey, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return errors.New("session not found or expired")
	}

	// Шаг 2: проверяем владельца перед удалением
	val, err := s.rdb.Get(ctx, seatKey).Bytes()
	if err != nil {
		return errors.New("seat data not found")
	}

	var b Booking
	if err := json.Unmarshal(val, &b); err != nil {
		return err
	}

	if b.UserID != userID {
		return errors.New("unauthorized")
	}

	// Шаг 3: удаляем оба ключа атомарно через pipeline
	pipe := s.rdb.Pipeline()
	pipe.Del(ctx, seatKey)
	pipe.Del(ctx, sessionKey(sessionID))
	_, err = pipe.Exec(ctx)

	return err
}

func (s *RedisStore) ConfirmSeat(ctx context.Context, sessionID, userID string) (Booking, error) {
	seatKey, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return Booking{}, errors.New("session not found or expired")
	}

	val, err := s.rdb.Get(ctx, seatKey).Bytes()
	if err != nil {
		return Booking{}, errors.New("seat data not found")
	}

	var b Booking
	if err := json.Unmarshal(val, &b); err != nil {
		return Booking{}, err
	}

	if b.UserID != userID {
		return Booking{}, errors.New("unauthorized")
	}

	b.Status = "confirmed"
	updated, _ := json.Marshal(b)

	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, seatKey, updated, 0) // 0 = без TTL, место занято навсегда
	pipe.Del(ctx, sessionKey(sessionID))
	_, err = pipe.Exec(ctx)

	return b, err
}
