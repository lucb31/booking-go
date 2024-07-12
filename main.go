package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
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
		}}
}

func main() {
	// Setup logger
	logger := log.Default()
	data := initData()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		jwt, err := c.Cookie("Jwt-Token")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			// c.HTML(http.StatusUnauthorized, "login.html", LoginResponse{jwt, err.Error()})
			return
		}
		token, err := VerifyJWT(jwt)
		if err != nil {
			//c.HTML(http.StatusUnauthorized, "login.html", LoginResponse{jwt, err.Error()})
			c.Redirect(http.StatusFound, "/login")
			return
		}
		logger.Print("Authorized request to index", token)
		c.HTML(http.StatusOK, "index.html", data)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", LoginResponse{"", ""})
	})
	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("Jwt-Token", "", 0, "", "localhost", true, true)
		c.Redirect(http.StatusFound, "/login")
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
	r.GET("/todos", func(c *gin.Context) {
		// Add a todo to demonstrate data change
		data.Todos = append(data.Todos, Todo{Title: "just another! ONE MORE", Done: false})
		// Render 'todos' template block
		c.HTML(http.StatusOK, "todos", data)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Run("0.0.0.0:8000")
}
