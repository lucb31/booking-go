package booking

import (
	"encoding/json"
	"time"
)

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
