package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)
	

	type Todo struct {
		ID        int    `json:"id"`
		Completed bool   `json:"completed"`
		Body      string `json:"body"`
	}
	
	var db *sql.DB
	
	func main() {
		fmt.Println("hello world")
		app := fiber.New()
	
		if os.Getenv("ENV") != "production" {
			// Load the .env file if not in production
			err := godotenv.Load(".env")
			if err != nil {
				log.Fatal("Error loading .env file:", err)
			}
		}
	
		DB_USER := os.Getenv("DB_USER")
		DB_PASSWORD := os.Getenv("DB_PASSWORD")
		DB_HOST := os.Getenv("DB_HOST")
		DB_PORT := os.Getenv("DB_PORT")
		DB_NAME := os.Getenv("DB_NAME")
	
		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
	
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
	
		fmt.Println("Connected to MySQL")
	
			
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	log.Fatal(app.Listen(":4000"))
}


func getTodos(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, completed, body FROM todos")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Completed, &todo.Body); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		todos = append(todos, todo)
	}
	return c.Status(200).JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := &Todo{}
	if err := c.BodyParser(todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	result, err := db.Exec("INSERT INTO todos (completed, body) VALUES (?, ?)", todo.Completed, todo.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	todo.ID = int(id)
	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := db.Exec("UPDATE todos SET completed = ? WHERE id = ?", true, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
