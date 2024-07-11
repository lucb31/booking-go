package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
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
	data := initData()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", data)
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
