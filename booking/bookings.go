package booking

import (
	"fmt"
	"time"
)

// Converts date in format "yyyy-mm-dd" & time in format "hh:mm" into unix timestamp
func TimeFromDateAndTime(dateString string, timeString string) (time.Time, error) {
	s := fmt.Sprintf("%s %s", dateString, timeString)
	res, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		return time.Time{}, err
	}
	return res, nil
}
