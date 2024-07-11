package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
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

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := TodoPageData{
		PageTitle: "My list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: false},
			{Title: "Task 3", Done: true},
		}}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data)
	})
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		data.Todos = append(data.Todos, Todo{Title: "just another! ONE MORE", Done: false})
		tmpl.ExecuteTemplate(w, "todos", data)
	})
	http.ListenAndServe(":8000", nil)
	return
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Run("0.0.0.0:8000")
}
