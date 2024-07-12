package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
	Rooms     []Room
}

type Room struct {
	Id    int
	Title string
}

type LoginResponse struct {
	Jwt          string
	ErrorMessage string
}

func initData() TodoPageData {
	return TodoPageData{
		PageTitle: "My list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: false},
			{Title: "Task 3", Done: true},
		},
		Rooms: []Room{
			{Id: 1, Title: "Room A"},
			{Id: 2, Title: "Room B"},
			{Id: 3, Title: "Room C"},
		}}
}

func main() {
	// Setup logger
	logger := log.Default()
	data := initData()

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
			c.HTML(http.StatusOK, "index.html", data)
		})
		authenticated.GET("/logout", func(c *gin.Context) {
			c.SetCookie("Jwt-Token", "", 0, "", "localhost", true, true)
			c.Redirect(http.StatusFound, "/login")
		})
		authenticated.GET("/todos", func(c *gin.Context) {
			// Add a todo to demonstrate data change
			data.Todos = append(data.Todos, Todo{Title: "just another! ONE MORE", Done: false})
			// Render 'todos' template block
			c.HTML(http.StatusOK, "todos", data)
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
				idx := slices.IndexFunc(data.Rooms, func(room Room) bool { return room.Id == idParam })
				if idx == -1 {
					c.HTML(http.StatusBadRequest, "", "")
					return
				}
				logger.Print("Trying to remove el", idx)
				// Remove element by index
				data.Rooms = append(data.Rooms[:idx], data.Rooms[idx+1:]...)
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
				data.Rooms = append(data.Rooms, Room{rand.Int(), title})
				c.HTML(http.StatusOK, "rooms", data)
			})
		}
	}

	r.Run("0.0.0.0:8000")
}
