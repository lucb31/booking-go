package booking

import (
	"testing"
	"time"
)

func TestDateTimeToTimestamp_ReturnsCorrectUnixTimestamp(t *testing.T) {
	dateString := "2024-07-03"
	timeString := "20:54"
	res, err := TimeFromDateAndTime(dateString, timeString)
	if err != nil {
		t.Fatalf(`Unable to parse string "%q" or "%q": %v"`, dateString, timeString, err)
	}
	expectedTt := 1720040040

	unixTt := res.Unix()
	if int(unixTt) != expectedTt {
		t.Fatalf(`Expected "%d", but received "%d"`, expectedTt, unixTt)
	}
}

func TestAddBooking_FailsWithEndDateBeforeStartDate(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.Add(-1)

	_, err := AddBooking(1, 1, startDate, endDate)

	if err == nil {
		t.Fatalf("Expected booking creation to fail with end date < start date")
	}
}

func TestFilterBookings_ReturnsIntersectingBookings(t *testing.T) {
	layout := "2006-01-02 15:04"
	// Add booking test data
	startDate, _ := time.Parse(layout, "2024-01-03 8:00")
	endDate, _ := time.Parse(layout, "2024-01-03 15:00")
	AddBooking(1, 1, startDate, endDate)

	// Test: Filter start date > start Date & fitler end date > end date
	// Fully enclosed
	filterStartDate, _ := time.Parse(layout, "2024-01-03 6:00")
	filterEndDate, _ := time.Parse(layout, "2024-01-03 18:00")
	res := FindBookingsWithinTimeInterval(&filterStartDate, &filterEndDate)

	expectedLength := 1
	resLength := len(res)
	if resLength != expectedLength {
		t.Fatalf("Expected filter result length %d, received %d", expectedLength, resLength)
	}
	if res[0] != Bookings[0] {
		t.Fatalf("Expected booking '%v' to be included", Bookings[0])
	}
}

func TestFilterBookings_OutsideWorkingHours(t *testing.T) {
	layout := "2006-01-02 15:04"
	// Add booking test data
	startDate, _ := time.Parse(layout, "2024-07-08 23:00")
	endDate, _ := time.Parse(layout, "2024-07-09 19:00")
	AddBooking(1, 1, startDate, endDate)
	startDate, _ = time.Parse(layout, "2024-07-08 01:00")
	endDate, _ = time.Parse(layout, "2024-07-09 07:00")
	bIncluded, _ := AddBooking(1, 1, startDate, endDate)

	// Test: Filter start date > start Date & fitler end date > end date
	// Fully enclosed
	filterStartDate, _ := time.Parse(layout, "2024-07-08 8:00")
	filterEndDate, _ := time.Parse(layout, "2024-07-08 17:00")
	res := FindBookingsWithinTimeInterval(&filterStartDate, &filterEndDate)

	expectedLength := 1
	resLength := len(res)
	if resLength != expectedLength {
		t.Fatalf("Expected filter result length %d, received %d", expectedLength, resLength)
	}
	if res[0] != *bIncluded {
		t.Fatalf("Expected booking '%v' to be included", Bookings[0])
	}
}
