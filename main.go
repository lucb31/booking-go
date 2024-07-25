package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"lucb31/booking-go/booking"
	"lucb31/booking-go/calendar"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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

var logger = log.Default()
var bookingRepo booking.BookingRepository
var userRepo booking.UserRepository
var roomRepo booking.RoomsRepository

func main() {
	// Initialize router
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets/")

	// Initialize DB
	db, err := sqlx.Connect("sqlite3", "file:test.db")
	if err != nil {
		log.Fatalln(err)
	}
	// Init repos
	userRepo = booking.NewUserRepositorySQLite(db)
	if err := userRepo.Migrate(); err != nil {
		log.Fatalln(err)
	}
	roomRepo = booking.NewRoomsRepositorySQLite(db)
	if err := roomRepo.Migrate(); err != nil {
		log.Fatalln(err)
	}
	bookingRepo = booking.NewBookingRepositorySQLite(db, userRepo, roomRepo)
	if err = bookingRepo.Migrate(); err != nil {
		log.Fatalln(err)
	}
	// Seed test data
	if err := userRepo.SeedTestData(); err != nil {
		log.Fatalln(err)
	}
	if err := roomRepo.SeedTestData(); err != nil {
		log.Fatalln(err)
	}
	if err := bookingRepo.SeedTestData(); err != nil {
		log.Fatalln(err)
	}

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
			data, err := getBookingPageData()
			if err != nil {
				logger.Panic(err)
				c.HTML(http.StatusOK, "index.html", BookingPageData{})
				return
			}
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
	roomNumericId, err := strconv.ParseInt(roomId, 10, 64)
	if err != nil {
		return err
	}
	userNumericId, err := strconv.ParseInt(userId, 10, 64)
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

	_, err = bookingRepo.Create(booking.Booking{Room: booking.Room{Id: roomNumericId}, User: booking.User{Id: userNumericId}, StartTime: startAt, EndTime: endAt})
	if err != nil {
		return err
	}
	data, err := getBookingPageData()
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "bookings", data)
	return nil
}

func handleDeleteBookingRequest(c *gin.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	err = bookingRepo.Delete(id)
	if err != nil {
		return err
	}
	data, err := getBookingPageData()
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "bookings", data)
	return nil
}

func getBookingPageData() (BookingPageData, error) {
	bookings, err := bookingRepo.GetAll()
	if err != nil {
		return BookingPageData{Error: err.Error()}, err
	}
	rooms, err := roomRepo.GetAll()
	if err != nil {
		return BookingPageData{Error: err.Error()}, err
	}
	users, err := userRepo.GetAll()
	if err != nil {
		return BookingPageData{Error: err.Error()}, err
	}
	return BookingPageData{pointerSliceToValueSlice(bookings), pointerSliceToValueSlice(rooms), pointerSliceToValueSlice(users), ""}, nil
}

func pointerSliceToValueSlice[t comparable](vals []*t) []t {
	res := make([]t, len(vals))
	for idx, val := range vals {
		res[idx] = *val
	}
	return res
}

func handleEditBookingRequest(c *gin.Context) error {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	record, err := bookingRepo.GetById(idParam)
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "booking-modal", BookingDetailData{Booking: *record})
	return nil
}

func handleUpdateBookingRequest(c *gin.Context) error {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	titleParam := c.Request.FormValue("title")
	record, err := bookingRepo.GetById(idParam)
	if err != nil {
		return err
	}
	record.Title = titleParam
	// Todo Seems to be triggering update twice. Need to investigate
	c.Header("HX-Trigger", "calendar-update")
	c.HTML(http.StatusOK, "booking-modal-form", BookingDetailData{Booking: *record})
	return nil
}

func handleDeleteRoomRequest(c *gin.Context) {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: err.Error()})
		return
	}
	err = roomRepo.Delete(idParam)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: err.Error()})
		return
	}

	rooms, err := roomRepo.GetAll()
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: err.Error()})
	}
	data := RoomPageData{pointerSliceToValueSlice(rooms)}
	c.HTML(http.StatusOK, "rooms", data)
}

func handleAddRoomRequest(c *gin.Context) {
	title, exists := c.GetPostForm("title")
	if !exists || len(title) == 0 {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: "Title cannot be empty"})
		return
	}
	_, err := roomRepo.Create(booking.Room{Title: title})
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: err.Error()})
		return
	}
	rooms, err := roomRepo.GetAll()
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "rooms", BookingPageData{Error: err.Error()})
	}
	data := RoomPageData{pointerSliceToValueSlice(rooms)}
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
	var service calendar.CalendarService = calendar.NewService(bookingRepo)
	dayData, err := service.GetCalendarDayData(year, week)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "calendar.html", CalendarData{})
	}
	data := CalendarData{service.GenerateTimeMarkers(), dayData, week, nextWeek, week - 1}
	c.HTML(http.StatusOK, "calendar.html", data)
}
