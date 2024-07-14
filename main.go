package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lucb31/booking-go/booking"
	"lucb31/booking-go/calendar"
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
	r.POST("/login", func(c *gin.Context) {
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
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

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
			// DELETE room by ID
			roomEndpoints.DELETE("/:id", func(c *gin.Context) {
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
			})

			// ADD room with title
			roomEndpoints.POST("/", func(c *gin.Context) {
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
			})
		}
		bookingEndpoints := authenticated.Group("/bookings")
		{
			// DELETE booking by ID
			bookingEndpoints.DELETE("/:id", func(c *gin.Context) {
				err := booking.RemoveBookingByIdString(c.Param("id"))
				if err != nil {
					c.HTML(http.StatusUnprocessableEntity, "bookings", newErrBokingPageData(err))
					return
				}
				data := getBookingPageData()
				c.HTML(http.StatusOK, "bookings", data)
			})

			// ADD booking
			bookingEndpoints.POST("/", func(c *gin.Context) {
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
			})
		}
		authenticated.GET("/calendar", func(c *gin.Context) {
			data := CalendarData{calendar.GenerateTimeMarkers(), calendar.GetCalendarDayData()}
			c.HTML(http.StatusOK, "calendar.html", data)
		})
	}

	r.Run("0.0.0.0:8000")
}
