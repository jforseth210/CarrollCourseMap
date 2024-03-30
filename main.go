package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	println("Starting...")
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")

	db, err := gorm.Open(sqlite.Open("classes.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	router.GET("/", func(c *gin.Context) {
		var buildings []Building
		var classes []Class
		db.Find(&buildings)
		db.Preload("Room.Building").Find(&classes)
		c.HTML(http.StatusOK, "index.html", gin.H{"Buildings": buildings, "Classes": classes})
	})
	db.AutoMigrate(&Building{}, &Room{}, &Class{})
	/*
		building := Building{Name: "Simperman", Latitude: 46.599809, Longitude: -112.037425}
		building1 := Building{Name: "O'Connell", Latitude: 46.600769, Longitude: -112.04027}
		building2 := Building{Name: "St. Charles", Latitude: 46.600854, Longitude: -112.038613}
		building3 := Building{Name: "All Saints' Chapel", Latitude: 46.601338, Longitude: -112.038511}
		building4 := Building{Name: "Library", Latitude: 46.601888, Longitude: -112.0381}
		building5 := Building{Name: "St. Albert", Latitude: 46.600557, Longitude: -112.037722}
		building6 := Building{Name: "Borromeo", Latitude: 46.599731, Longitude: -112.03917}
		building7 := Building{Name: "Cube", Latitude: 46.599888, Longitude: -112.040354}
		building8 := Building{Name: "Civil Engineering", Latitude: 46.599567, Longitude: -112.036194}
		building9 := Building{Name: "Canine Center", Latitude: 46.599641, Longitude: -112.035523}

		db.Create(&building)  // pass pointer of data to Create
		db.Create(&building1) // pass pointer of data to Create
		db.Create(&building2) // pass pointer of data to Create
		db.Create(&building3) // pass pointer of data to Create
		db.Create(&building4) // pass pointer of data to Create
		db.Create(&building5) // pass pointer of data to Create
		db.Create(&building6) // pass pointer of data to Create
		db.Create(&building7) // pass pointer of data to Create
		db.Create(&building8) // pass pointer of data to Create
		db.Create(&building9) // pass pointer of data to Create
	*/
	router.GET("/api/rooms/:id", func(c *gin.Context) {
		id := c.Param("id")
		var rooms []Room
		db.Find(&rooms, "building_id =?", id)
		c.JSON(http.StatusOK, rooms)
	})

	router.POST("/api/rooms/add-room", func(c *gin.Context) {
		var room Room
		if err := c.ShouldBindJSON(&room); err != nil {
			c.JSON(400, gin.H{"error": "Invalid data"})
			return
		}
		// Save room to the database
		db.Create(&room)
		// Return the room ID
		c.JSON(200, gin.H{"roomId": room.ID})
	})

	router.POST("/api/classes/add-class", func(c *gin.Context) {
		// Parse form data
		className := c.PostForm("name")
		courseCode := c.PostForm("code")
		professorName := c.PostForm("professor")
		roomID, _ := strconv.Atoi(c.PostForm("roomId"))
		var room Room
		db.First(&room, roomID)

		// Create a new class
		newClass := Class{
			Name:      className,
			Code:      courseCode, // Set the description if needed
			RoomID:    uint(roomID),
			Room:      room,
			Professor: professorName,
		}

		// Save both class and professor to the database
		db.Create(&newClass)

		// Redirect or send a success response
		http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
	})

	// Migrate the schema
	println("Migrating schema...")
	println("Running server...")
	router.Run(":8000")
}
