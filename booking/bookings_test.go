package booking

import (
	"testing"
)

func TestDateTimeToTimestamp_ReturnsCorrectUnixTimestamp(t *testing.T) {
	dateString := "2024-07-03"
	timeString := "20:54"
	timestamp, err := DateTimeToTimestamp(dateString, timeString)
	if err != nil {
		t.Fatalf(`Unable to parse string "%q" or "%q": %v"`, dateString, timeString, err)
	}
	expectedTt := 1720040040
	if timestamp != expectedTt {
		t.Fatalf(`Expected "%d", but received "%d"`, expectedTt, timestamp)
	}

}

func TestAddBooking_FailsWithEndDateBeforeStartDate(t *testing.T) {
	var startDate int = 5000
	var endDate int = 3000

	_, err := AddBooking(1, 1, startDate, endDate)

	if err == nil {
		t.Fatalf("Expected booking creation to fail with end date < start date")
	}
}
