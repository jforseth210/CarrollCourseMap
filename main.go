package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	println("Starting...")
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	db, err := gorm.Open(sqlite.Open("classes.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	println("Migrating schema...")
	db.AutoMigrate(&Building{}, &Room{}, &Class{}, &Professor{}, &ClassPeriod{})
	println("Running server...")
	router.Run(":8000")
}
