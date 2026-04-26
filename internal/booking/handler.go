package booking

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc *Service
}

func NewHandler(svc *Service) *handler {
	return &handler{svc}
}

type holdSeatRequest struct {
	UserID string `json:"user_id"`
}

// HoldSeat: POST /movies/:movieID/seats/:seatID/hold
func (h *handler) HoldSeat(c *gin.Context) {
	movieID := c.Param("movieID")
	seatID := c.Param("seatID")

	var req holdSeatRequest
	// ShouldBindJSON заменяет json.NewDecoder — плюс автоматически
	// отвечает 400, если тело невалидно
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.svc.Book(Booking{
		UserID:  req.UserID,
		SeatID:  seatID,
		MovieID: movieID,
	})
	if err != nil {
		// ErrSeatAlreadyBooked — это ожидаемая бизнес-ошибка, не 500
		if err == ErrSeatAlreadyBooked {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session_id": session.ID,
		"movie_id":   session.MovieID,
		"seat_id":    session.SeatID,
		"expires_at": session.ExpiresAt.Format(time.RFC3339),
	})
}

// ListSeats: GET /movies/:movieID/seats
func (h *handler) ListSeats(c *gin.Context) {
	movieID := c.Param("movieID")
	bookings := h.svc.ListBookings(movieID)

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID:    b.SeatID,
			UserID:    b.UserID,
			Booked:    true,
			Confirmed: b.Status == "confirmed",
		})
	}

	c.JSON(http.StatusOK, seats)
}

type seatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}

// ConfirmSession: PUT /sessions/:sessionID/confirm
func (h *handler) ConfirmSession(c *gin.Context) {
	sessionID := c.Param("sessionID")

	var req holdSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	session, err := h.svc.ConfirmSeat(c.Request.Context(), sessionID, req.UserID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.ID,
		"movie_id":   session.MovieID,
		"seat_id":    session.SeatID,
		"user_id":    req.UserID,
		"status":     session.Status,
	})
}

// ReleaseSession: DELETE /sessions/:sessionID
func (h *handler) ReleaseSession(c *gin.Context) {
	sessionID := c.Param("sessionID")

	var req holdSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if err := h.svc.ReleaseSeat(c.Request.Context(), sessionID, req.UserID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 204 No Content — стандартный ответ при успешном удалении без тела
	c.Status(http.StatusNoContent)
}

type Movie struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

// GET /movies
func (h *handler) ListMovies(c *gin.Context) {
	movies := []Movie{
		{ID: "screen-1", Title: "Dune", Rows: 4, SeatsPerRow: 5},
		{ID: "screen-2", Title: "Richie Rich", Rows: 2, SeatsPerRow: 10},
	}
	c.JSON(http.StatusOK, movies)
}
