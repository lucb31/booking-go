package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lucb31/booking-go/booking"
)

type RoomPageData struct {
	Rooms []booking.Room
}

type LoginResponse struct {
	Jwt          string
	ErrorMessage string
}

func main() {
	// Setup logger
	logger := log.Default()

	// Initialize router
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets/")

	// Unauthorized routes
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", LoginResponse{"", ""})
	})
	r.POST("/login", func(c *gin.Context) {
		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")
		logger.Print("Received", username, password)

		jwt, err := LoginRequest(username, password)
		if err != nil {
			logger.Print("Failed login request: ", err)
			c.HTML(http.StatusUnauthorized, "login.html", LoginResponse{jwt, err.Error()})
			return
		}
		logger.Print("Received", jwt, err)

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
			data := RoomPageData{booking.Rooms}
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
					c.HTML(http.StatusBadRequest, "", "")
					return
				}
				err = booking.RemoveById(idParam)
				if err != nil {
					c.HTML(http.StatusBadRequest, "", "")
					return
				}

				data := RoomPageData{booking.Rooms}
				c.HTML(http.StatusOK, "rooms", data)
			})

			// ADD room with title
			roomEndpoints.POST("/", func(c *gin.Context) {
				title, exists := c.GetPostForm("title")
				if !exists {
					c.HTML(http.StatusBadRequest, "", "")
					logger.Print("Title cannot be empty")
					return
				}
				_, err := booking.Add(title)
				if err != nil {
					c.HTML(http.StatusBadRequest, "", "")
					logger.Print("Unable to create room: ", err)
					return
				}
				data := RoomPageData{booking.Rooms}
				c.HTML(http.StatusOK, "rooms", data)
			})
		}
	}

	r.Run("0.0.0.0:8000")
}
