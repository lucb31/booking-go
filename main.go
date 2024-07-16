package main

import (
	"errors"
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

func getBookingPageData() BookingPageData {
	return BookingPageData{booking.Bookings, booking.Rooms, booking.Users, ""}
}

func newErrBokingPageData(err error) BookingPageData {
	logger.Print("Failed booking page data request: ", err)
	return BookingPageData{booking.Bookings, booking.Rooms, booking.Users, err.Error()}
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
			data := getBookingPageData()
			c.HTML(http.StatusOK, "index.html", data)
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
			bookingEndpoints.POST("/", handleAddBookingRequest)
			bookingEndpoints.DELETE("/:id", handleDeleteBookingRequest)
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

func handleAddBookingRequest(c *gin.Context) {
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
		c.HTML(http.StatusUnprocessableEntity, "bookings", newErrBokingPageData(err))
		return
	}
	userNumericId, err := strconv.Atoi(userId)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "bookings", newErrBokingPageData(err))
		return
	}

	// Convert date inputs into unix TT
	startAt, err := booking.TimeFromDateAndTime(startDate, startTime)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "bookings", newErrBokingPageData(err))
		return
	}
	endAt, err := booking.TimeFromDateAndTime(endDate, endTime)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "bookings", newErrBokingPageData(err))
		return
	}

	_, err = booking.AddBooking(roomNumericId, userNumericId, startAt, endAt)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "bookings", newErrBokingPageData(err))
		return
	}
	data := getBookingPageData()
	c.HTML(http.StatusOK, "bookings", data)
}

func handleDeleteBookingRequest(c *gin.Context) {
	err := booking.RemoveBookingByIdString(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "bookings", newErrBokingPageData(err))
		return
	}
	data := getBookingPageData()
	c.HTML(http.StatusOK, "bookings", data)
}

func handleDeleteRoomRequest(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", newErrBokingPageData(err))
		return
	}
	err = booking.RemoveRoomById(idParam)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", newErrBokingPageData(err))
		return
	}

	data := RoomPageData{booking.Rooms}
	c.HTML(http.StatusOK, "rooms", data)
}

func handleAddRoomRequest(c *gin.Context) {
	title, exists := c.GetPostForm("title")
	if !exists || len(title) == 0 {
		c.HTML(http.StatusUnprocessableEntity, "rooms", newErrBokingPageData(errors.New("Title cannot be empty")))
		return
	}
	_, err := booking.AddRoom(title)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", newErrBokingPageData(err))
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
