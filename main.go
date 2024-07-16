package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"lucb31/booking-go/booking"
	"lucb31/booking-go/calendar"

	"github.com/gin-gonic/gin"
)

type RoomPageData struct {
	Rooms []booking.Room
}

type BookingPageData struct {
	Bookings []booking.Booking
	Rooms    []booking.Room
	Users    []booking.User
	Error    string
}

type BookingDetailData struct {
	Booking booking.Booking
	Error   string
}

type LoginResponse struct {
	Jwt          string
	ErrorMessage string
}

type CalendarData struct {
	TimeMarkers []string
	DayData     []calendar.CalendarDayData
	Cw          int
	NextCw      int
	PrevCw      int
}

// Setup logger
var logger = log.Default()

func main() {
	// Initialize router
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets/")

	// Seed test data
	booking.InitTestBookings()

	// Unauthorized routes
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", LoginResponse{"", ""})
	})
	r.POST("/login", handleLoginRequest)

	// Authorized routes
	authenticated := r.Group("/")
	authenticated.Use(AuthMiddleware())
	{
		authenticated.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", BookingPageData{booking.Bookings, booking.Rooms, booking.Users, ""})
		})
		authenticated.GET("/logout", func(c *gin.Context) {
			c.SetCookie("Jwt-Token", "", 0, "", "localhost", true, true)
			c.Redirect(http.StatusFound, "/login")
		})

		roomEndpoints := authenticated.Group("/rooms")
		{
			roomEndpoints.DELETE("/:id", handleDeleteRoomRequest)
			roomEndpoints.POST("/", handleAddRoomRequest)
		}
		bookingEndpoints := authenticated.Group("/bookings")
		{
			bookingEndpoints.POST("/", makeBookingRequest(handleAddBookingRequest))
			bookingEndpoints.GET("/:id", makeBookingRequest(handleEditBookingRequest))
			bookingEndpoints.DELETE("/:id", makeBookingRequest(handleDeleteBookingRequest))
			bookingEndpoints.PATCH("/:id", makeBookingModalRequest(handleUpdateBookingRequest))
		}
		authenticated.GET("/calendar", handleGetCalendarRequest)
	}

	r.Run("0.0.0.0:8000")
}

func handleLoginRequest(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	logger.Printf("Login request for user %s", username)

	jwt, err := LoginRequest(username, password)
	if err != nil {
		logger.Print("Failed login request: ", err)
		c.HTML(http.StatusUnauthorized, "login.html", LoginResponse{jwt, err.Error()})
		return
	}
	c.SetCookie("Jwt-Token", jwt, 86400, "", "localhost", true, true)
	c.Redirect(http.StatusFound, "/")
}

// Middleware for booking request errors
func makeBookingRequest(h func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := h(c)
		if err != nil {
			c.HTML(http.StatusUnprocessableEntity, "bookings", BookingPageData{Error: err.Error()})
			return
		}
	}
}

// Middleware for booking-modal request errors
func makeBookingModalRequest(h func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := h(c)
		if err != nil {
			c.HTML(http.StatusUnprocessableEntity, "booking-modal", BookingDetailData{Error: err.Error()})
			return
		}
	}
}

func handleAddBookingRequest(c *gin.Context) error {
	// Fetch inputs
	roomId, _ := c.GetPostForm("roomId")
	userId, _ := c.GetPostForm("userId")
	startDate, _ := c.GetPostForm("startDate")
	startTime, _ := c.GetPostForm("startTime")
	endDate, _ := c.GetPostForm("endDate")
	endTime, _ := c.GetPostForm("endTime")

	// Convert string ids to numeric
	roomNumericId, err := strconv.Atoi(roomId)
	if err != nil {
		return err
	}
	userNumericId, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}

	// Convert date inputs into unix TT
	startAt, err := booking.TimeFromDateAndTime(startDate, startTime)
	if err != nil {
		return err
	}
	endAt, err := booking.TimeFromDateAndTime(endDate, endTime)
	if err != nil {
		return err
	}

	_, err = booking.AddBooking(roomNumericId, userNumericId, startAt, endAt)
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "bookings", BookingPageData{booking.Bookings, booking.Rooms, booking.Users, ""})
	return nil
}

func handleDeleteBookingRequest(c *gin.Context) error {
	err := booking.RemoveBookingByIdString(c.Param("id"))
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "bookings", BookingPageData{booking.Bookings, booking.Rooms, booking.Users, ""})
	return nil
}

func handleEditBookingRequest(c *gin.Context) error {
	idParam := c.Param("id")
	record := booking.FindBookingByIdString(idParam)
	if record == nil {
		return errors.New(fmt.Sprintf("Could not find booking for id %s", idParam))
	}
	c.HTML(http.StatusOK, "booking-modal", BookingDetailData{Booking: *record})
	return nil
}

func handleUpdateBookingRequest(c *gin.Context) error {
	idParam := c.Param("id")
	titleParam := c.Request.FormValue("title")
	record := booking.FindBookingByIdString(idParam)
	if record == nil {
		return errors.New(fmt.Sprintf("Could not find booking for id %s", idParam))
	}
	record.Title = titleParam
	// Todo Seems to be triggering update twice. Need to investigate
	c.Header("HX-Trigger", "calendar-update")
	c.HTML(http.StatusOK, "booking-modal-form", BookingDetailData{Booking: *record})
	return nil
}

func handleDeleteRoomRequest(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: err.Error()})
		return
	}
	err = booking.RemoveRoomById(idParam)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: err.Error()})
		return
	}

	data := RoomPageData{booking.Rooms}
	c.HTML(http.StatusOK, "rooms", data)
}

func handleAddRoomRequest(c *gin.Context) {
	title, exists := c.GetPostForm("title")
	if !exists || len(title) == 0 {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: "Title cannot be empty"})
		return
	}
	_, err := booking.AddRoom(title)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: err.Error()})
		return
	}
	data := RoomPageData{booking.Rooms}
	c.HTML(http.StatusOK, "rooms", data)
}

func handleGetCalendarRequest(c *gin.Context) {
	week, err := strconv.Atoi(c.Query("week"))
	// Fallback to current week if invalid or none provided
	if err != nil || week < 1 || week > 53 {
		_, week = time.Now().ISOWeek()
	}
	year, err := strconv.Atoi(c.Query("year"))
	// Fallback to current year if invalid or none provided
	if err != nil || year < 1 || year > 1000 {
		year, _ = time.Now().ISOWeek()
	}
	// Disable next week button if last week reached
	nextWeek := week + 1
	if nextWeek > 53 {
		nextWeek = 0
	}
	data := CalendarData{calendar.GenerateTimeMarkers(), calendar.GetCalendarDayData(year, week), week, nextWeek, week - 1}
	c.HTML(http.StatusOK, "calendar.html", data)
}
