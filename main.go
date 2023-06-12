package main

import (
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	gorm.Model
	Task   string
	Status string
}

func main() {
	dsn := "root:password@tcp(localhost:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Todo{})

	r := gin.Default()

	// GET /todos
	r.GET("/todos", func(c *gin.Context) {
		var todos []Todo
		result := db.Find(&todos)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, todos)
	})

	// POST /todos
	r.POST("/todos", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := db.Create(&todo)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusCreated, todo)
	})

	// GET /todos/:id
	r.GET("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")

		var todo Todo
		result := db.First(&todo, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		c.JSON(http.StatusOK, todo)
	})

	// PUT /todos/:id
	r.PUT("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")

		var todo Todo
		result := db.First(&todo, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result = db.Save(&todo)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)
	})

	// DELETE /todos/:id
	r.DELETE("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")

		var todo Todo
		result := db.First(&todo, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		result = db.Delete(&todo)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
	})

	if err := r.Run(":8089"); err != nil {
		log.Fatal(err)
	}
}
