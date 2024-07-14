package booking

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"
)

type Booking struct {
	Id        int
	Room      *Room
	User      *User
	StartTime int
	EndTime   int
}

var Bookings []Booking

func RemoveBookingById(id int) error {
	idx := slices.IndexFunc(Bookings, func(booking Booking) bool { return booking.Id == id })
	if idx == -1 {
		return errors.New("Unknown booking id")
	}
	// Remove element by index
	Bookings = append(Bookings[:idx], Bookings[idx+1:]...)
	return nil
}

func AddBooking(roomId int, userId int, startAt int, endAt int) (*Booking, error) {
	// Find room
	room, err := GetRoomById(roomId)
	if err != nil {
		return nil, err
	}
	// Find user
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	// Validate dates
	if startAt > endAt {
		return nil, errors.New("Start date cannot be after end date")
	}

	newBooking := Booking{rand.IntN(20000), room, user, startAt, endAt}
	Bookings = append(Bookings, newBooking)
	return &newBooking, nil
}

// Converts date in format "yyyy-mm-dd" & time in format "hh:mm" into unix timestamp
func DateTimeToTimestamp(dateString string, timeString string) (int, error) {
	s := fmt.Sprintf("%s %s", dateString, timeString)
	res, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		return 0, err
	}
	return int(res.Unix()), nil
}
