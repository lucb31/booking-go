package calendar

import (
	"fmt"
	"lucb31/booking-go/booking"
	"time"
)

type CalendarEvent struct {
	// Relative to workingHourStart
	StartHour int
	// Relative to workingHourStart
	EndHour int
	Booking *booking.Booking
}

type CalendarDayData struct {
	DayNum    int
	DayString string
	Events    []CalendarEvent
	//Bookings []booking.Booking
}

// Working hours will be split in numTimeMarkers blocks within the calendar
const numTimeMarkers = 10
const workingHourStart = 8
const workingHoursEnd = 17

func GenerateTimeMarkers() []string {
	// Generate time markers
	var timeMarkers [numTimeMarkers]string
	for i := 0; i < numTimeMarkers; i++ {
		workingHour := i + workingHourStart
		// Convert 24h hours into AM/PM format
		amPm := "AM"
		if workingHour > 11 {
			amPm = "PM"
		}
		timeMarkers[i] = fmt.Sprintf("%d %s", workingHour, amPm)
	}
	return timeMarkers[:]
}

func GetCalendarDayData() []CalendarDayData {
	workingDays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}

	currentTime := time.Now().AddDate(0, 0, 0)
	currentWeekday := currentTime.Weekday()

	offsetMonday := time.Monday - currentWeekday
	dateOfFirstMonday := currentTime.AddDate(0, 0, int(offsetMonday)-7)

	var dayData [5]CalendarDayData
	for idx, workingDay := range workingDays {
		dayString := workingDay.String()[0:3]
		workingTime := dateOfFirstMonday.AddDate(0, 0, idx)
		filterStartDate := time.Date(workingTime.Year(), workingTime.Month(), workingTime.Day(), workingHourStart, 0, 0, 0, workingTime.Location())
		filterEndDate := time.Date(workingTime.Year(), workingTime.Month(), workingTime.Day(), workingHoursEnd+1, 0, 0, 0, workingTime.Location())
		dayNum := workingTime.Day()

		// Filter bookings by calendar date
		filteredBookings := booking.FindBookingsWithinTimeInterval(&filterStartDate, &filterEndDate)

		// Map bookings to Event data
		events := make([]CalendarEvent, len(filteredBookings))
		for idx, b := range filteredBookings {
			events[idx] = mapBookingToCalendarEvent(&b, &filterStartDate, &filterEndDate)
		}
		dayData[idx] = CalendarDayData{dayNum, dayString, events}
	}

	return dayData[:]
}

func mapBookingToCalendarEvent(b *booking.Booking, startLimit *time.Time, endLimit *time.Time) CalendarEvent {
	relativeStartHour := 1
	if !b.StartTime.Before(*startLimit) {
		// Offset by starting work hour, starting at 1; cannot be lower than 1
		relativeStartHour = max(1, min(numTimeMarkers, b.StartTime.Hour()-workingHourStart+1))
	}
	relativeEndHour := numTimeMarkers
	if !b.EndTime.After(*endLimit) {
		// Offset by starting work hour, starting at 1; cannot be lower than numTimeMarkers
		relativeEndHour = max(1, min(numTimeMarkers, b.EndTime.Hour()-workingHourStart+1))
	}
	return CalendarEvent{relativeStartHour, relativeEndHour, b}
}
