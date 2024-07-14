package booking

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	Id          uuid.UUID
	Room        *Room
	User        *User
	StartTime   time.Time
	EndTime     time.Time
	Title       string
	Description string
}

func (b *Booking) Duration() time.Duration {
	return b.EndTime.Sub(b.StartTime)
}

// Return true, if the given Time is within the booking Interval
func (b *Booking) Within(t *time.Time) bool {
	// !After & !Before required to treat equals case correctly
	return !b.StartTime.After(*t) && !b.EndTime.Before(*t)
}

func (b *Booking) String() string {
	jsonBytes, err := json.Marshal(b)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// Constructor
func NewBooking(roomId int, userId int, startAt time.Time, endAt time.Time) (*Booking, error) {
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
	if !endAt.After(startAt) {
		return nil, errors.New("End date needs to be after start date")
	}

	description := fmt.Sprintf("%s to %s", startAt, endAt)
	id := uuid.New()
	newBooking := Booking{id, room, user, startAt, endAt, id.String(), description}
	return &newBooking, nil
}
