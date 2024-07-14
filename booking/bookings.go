package booking

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"time"

	"github.com/google/uuid"
)

var Bookings []Booking

func InitTestBookings() {
	Bookings = nil
	layout := "2006-01-02 15:04"
	// Add booking test data
	minStartDate, _ := time.Parse(layout, "2024-07-08 00:00")
	for i := 0; i < 10; i++ {
		offset := int64(rand.Intn(24 * 5))
		duration := int64(rand.Intn(25))
		startDate := minStartDate.Add(time.Duration(time.Hour * time.Duration(offset)))
		endDate := startDate.Add(time.Duration(time.Hour * time.Duration(duration)))
		AddBooking(1, 1, startDate, endDate)
	}
}

func RemoveBookingByIdString(idString string) error {
	id, err := uuid.Parse(idString)
	idx := slices.IndexFunc(Bookings, func(booking Booking) bool { return booking.Id == id })
	if idx == -1 || err != nil {
		return errors.New("Unknown booking id")
	}
	// Remove element by index
	Bookings = append(Bookings[:idx], Bookings[idx+1:]...)
	return nil
}

func AddBooking(roomId int, userId int, startAt time.Time, endAt time.Time) (*Booking, error) {
	newBooking, err := NewBooking(roomId, userId, startAt, endAt)
	if err != nil {
		return nil, err
	}
	Bookings = append(Bookings, *newBooking)
	log.Print("Booking added", newBooking)
	return newBooking, nil
}

// Returns all bookings that intersect with the given time interval
func FindBookingsWithinTimeInterval(startAt *time.Time, endAt *time.Time) []Booking {
	return filterBookings(func(b *Booking) bool {
		// Booking intersects with either start or end Time
		if b.Within(startAt) || b.Within(endAt) {
			return true
		}
		// Booking fully included in time interval
		if !startAt.After(b.StartTime) && !endAt.Before(b.EndTime) {
			return true
		}
		return false
	})
}

func filterBookings(f func(b *Booking) bool) []Booking {
	res := make([]Booking, 0)
	for _, b := range Bookings {
		if f(&b) {
			res = append(res, b)
		}
	}
	return res
}

// Converts date in format "yyyy-mm-dd" & time in format "hh:mm" into unix timestamp
func TimeFromDateAndTime(dateString string, timeString string) (time.Time, error) {
	s := fmt.Sprintf("%s %s", dateString, timeString)
	res, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		return time.Time{}, err
	}
	return res, nil
}