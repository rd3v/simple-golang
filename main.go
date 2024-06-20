package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Todo struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Done int `json:"done"`
}

var todos = []Todo{
	{ID: 1, Name: "Belajar Golang", Done: 1},
	{ID: 2, Name: "Belajar ReactJS", Done: 0},
	{ID: 3, Name: "Belajar Laravel", Done: 1},
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, World!")
	})

	e.GET("/todos", getAll)
	e.GET("/todos/:id", getById)
	e.POST("/todos", createTodo)
	e.PUT("/todos/:id", updateTodo)
	e.DELETE("/todos/:id", deleteTodo)	

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func getAll(c echo.Context) error {
	return c.JSON(http.StatusOK, todos)
}

func getById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	for _, todo := range todos {
		if todo.ID == id {
			return c.JSON(http.StatusOK, todo)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "To-do item not found"})	
}

// Handler to create a new to-do item
func createTodo(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Set a new ID for the to-do item
	todo.ID = getNextID()
	todos = append(todos, *todo)
	return c.JSON(http.StatusCreated, todo)
}

// Handler to update an existing to-do item by ID
func updateTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	updatedTodo := new(Todo)
	if err := c.Bind(updatedTodo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	for i, todo := range todos {
		if todo.ID == id {
			// Update the to-do item
			todos[i].Name = updatedTodo.Name
			todos[i].Done = updatedTodo.Done
			return c.JSON(http.StatusOK, todos[i])
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "To-do item not found"})
}

// Handler to delete a to-do item by ID
func deleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	for i, todo := range todos {
		if todo.ID == id {
			// Delete the to-do item
			todos = append(todos[:i], todos[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "To-do item not found"})
}

// Helper function to get the next ID for a new to-do item
func getNextID() int {
	if len(todos) == 0 {
		return 1
	}
	return todos[len(todos)-1].ID + 1
}

// Simple implementation of an integer minimum
// Adapted from: https://gobyexample.com/testing-and-benchmarking
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
