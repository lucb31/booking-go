package booking

import (
// "testing"
// "time"
)

const layout = "2006-01-02 15:04"

//func resetBookings() {
//	Bookings = []Booking{}
//}
//
//func TestDateTimeToTimestamp_ReturnsCorrectUnixTimestamp(t *testing.T) {
//	resetBookings()
//	dateString := "2024-07-03"
//	timeString := "20:54"
//	res, err := TimeFromDateAndTime(dateString, timeString)
//	if err != nil {
//		t.Fatalf(`Unable to parse string "%q" or "%q": %v"`, dateString, timeString, err)
//	}
//	expectedTt := 1720040040
//
//	unixTt := res.Unix()
//	if int(unixTt) != expectedTt {
//		t.Fatalf(`Expected "%d", but received "%d"`, expectedTt, unixTt)
//	}
//}
//
//func TestFilterBookings_ReturnsIntersectingBookings(t *testing.T) {
//	resetBookings()
//	layout := "2006-01-02 15:04"
//	// Add booking test data
//	startDate, _ := time.Parse(layout, "2024-01-03 8:00")
//	endDate, _ := time.Parse(layout, "2024-01-03 15:00")
//	_, err := AddBooking(1, 1, startDate, endDate)
//	if err != nil {
//		t.Fatalf("Unable to create booking from '%s' to '%s'", startDate, endDate)
//	}
//
//	// Test: Filter start date > start Date & fitler end date > end date
//	// Fully enclosed
//	filterStartDate, _ := time.Parse(layout, "2024-01-03 6:00")
//	filterEndDate, _ := time.Parse(layout, "2024-01-03 18:00")
//	res := FindBookingsWithinTimeInterval(&filterStartDate, &filterEndDate)
//
//	expectedLength := 1
//	resLength := len(res)
//	if resLength != expectedLength {
//		t.Fatalf("Expected filter result length %d, received %d", expectedLength, resLength)
//	}
//	if res[0] != Bookings[0] {
//		t.Fatalf("Expected booking '%v' to be included", Bookings[0])
//	}
//}
//
//func TestFilterBookings_OutsideWorkingHours(t *testing.T) {
//	resetBookings()
//	// Add booking starting AFTER end date
//	startDate, _ := time.Parse(layout, "2024-07-08 23:00")
//	endDate, _ := time.Parse(layout, "2024-07-09 19:00")
//	b, err := AddBooking(1, 1, startDate, endDate)
//	if err != nil {
//		t.Errorf("Unable to create booking from '%s' to '%s': %s", startDate, endDate, err)
//	}
//
//	filterStartDate, _ := time.Parse(layout, "2024-07-08 8:00")
//	filterEndDate, _ := time.Parse(layout, "2024-07-08 17:00")
//	res := FindBookingsWithinTimeInterval(&filterStartDate, &filterEndDate)
//
//	if len(res) != 0 {
//		t.Fatalf("Expected booking '%s' to not be included in interval '%s' - '%s'", b, filterStartDate, filterEndDate)
//	}
//}
//
//func TestFilterBookings_StartingBeforeStartEndingAfterEnd(t *testing.T) {
//	resetBookings()
//	// Add booking starting BEFORE filter start date & ending AFTER filter end date
//	startDate, _ := time.Parse(layout, "2024-07-08 01:00")
//	endDate, _ := time.Parse(layout, "2024-07-09 07:00")
//	b, err := AddBooking(1, 1, startDate, endDate)
//	if err != nil {
//		t.Errorf("Unable to create booking from '%s' to '%s': %s", startDate, endDate, err)
//	}
//
//	filterStartDate, _ := time.Parse(layout, "2024-07-08 8:00")
//	filterEndDate, _ := time.Parse(layout, "2024-07-08 17:00")
//	res := FindBookingsWithinTimeInterval(&filterStartDate, &filterEndDate)
//
//	if len(res) != 1 {
//		t.Fatalf("Expected booking '%s' to not be included in interval '%s' - '%s'", b, filterStartDate, filterEndDate)
//	}
//	if res[0] != *b {
//		t.Fatalf("Expected booking '%s' to be included in interval '%s' - '%s'", b, filterStartDate, filterEndDate)
//	}
//}
//
//func TestFilterBookings_StartingAfterStartEndingBeforeEnd(t *testing.T) {
//	resetBookings()
//	// Add booking starting AFTER filter start date & ending BEFORE filter end date
//	startDate, _ := time.Parse(layout, "2024-07-08 10:00")
//	endDate, _ := time.Parse(layout, "2024-07-09 22:00")
//	b, err := AddBooking(1, 1, startDate, endDate)
//	if err != nil {
//		t.Errorf("Unable to create booking from '%s' to '%s': %s", startDate, endDate, err)
//	}
//
//	filterStartDate, _ := time.Parse(layout, "2024-07-08 8:00")
//	filterEndDate, _ := time.Parse(layout, "2024-07-08 17:00")
//	res := FindBookingsWithinTimeInterval(&filterStartDate, &filterEndDate)
//
//	if len(res) != 1 {
//		t.Fatalf("Expected booking '%s' to not be included in interval '%s' - '%s'", b, filterStartDate, filterEndDate)
//	}
//	if res[0] != *b {
//		t.Fatalf("Expected booking '%s' to be included in interval '%s' - '%s'", b, filterStartDate, filterEndDate)
//	}
//}
//
//func TestAddBooking_FailsWhenConflicting(t *testing.T) {
//	resetBookings()
//	// Add booking test data
//	startDate, _ := time.Parse(layout, "2024-07-08 8:00")
//	endDate, _ := time.Parse(layout, "2024-07-09 17:00")
//	AddBooking(1, 1, startDate, endDate)
//
//	// Adding the same booking fails
//	_, err := AddBooking(1, 1, startDate, endDate)
//	if err == nil {
//		t.Fatalf("Expected adding duplicate booking to fail")
//	}
//
//	// Adding an intersecting booking fails
//	startDate, _ = time.Parse(layout, "2024-07-08 15:00")
//	endDate, _ = time.Parse(layout, "2024-07-09 22:00")
//	_, err = AddBooking(1, 1, startDate, endDate)
//	if err == nil {
//		t.Fatalf("Expected adding intersecting booking to fail")
//	}
//}
