package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Main application routing and logic
// Author: Justin Forseth

func main() {
	// Initialize router
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")

	// Initialize database
	db, err := gorm.Open(sqlite.Open("classes.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Building{}, &Room{}, &Class{})

	// Populate buildings if empty
	var building Building
	if db.Limit(1).Find(&building).RowsAffected == 0 {
		CreateCarrollBuildings(db)
	}
	router.Static("/static", "./static")
	// The main page
	router.GET("/", func(c *gin.Context) {
		// Get list of buildings for the selector
		var buildings []Building
		db.Find(&buildings)

		// Get list of classes (with their locations)
		var classes []Class
		db.Preload("Room.Building").Find(&classes)

		// Render the page
		c.HTML(http.StatusOK, "index.html", gin.H{"Buildings": buildings, "Classes": classes})
	})

	// Get a building's rooms
	router.GET("/api/rooms/:building_id", func(c *gin.Context) {
		// Get the id
		buildingID := c.Param("building_id")

		// Get the building's rooms
		var rooms []Room
		db.Find(&rooms, "building_id =?", buildingID)

		// Return the rooms (this assumes the client is trusted with ALL room data)
		c.JSON(http.StatusOK, rooms)
	})

	// Add a new room
	router.POST("/api/rooms/add-room", func(c *gin.Context) {
		// Attempt to create room
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

	// Add a new class
	router.POST("/api/classes/add-class", func(c *gin.Context) {
		// Parse for data
		className := c.PostForm("name")
		courseCode := c.PostForm("code")
		professorName := c.PostForm("professor")
		roomID, _ := strconv.Atoi(c.PostForm("roomId"))

		// Get the room by ID
		var room Room
		db.First(&room, roomID)

		// Create the class
		newClass := Class{
			Name:      className,
			Code:      courseCode,
			RoomID:    uint(roomID),
			Room:      room,
			Professor: professorName,
		}

		// Save class to the database
		db.Create(&newClass)

		// Redirect to main page
		http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
	})

	router.Run()
}

// Create building models for important building on campus
func CreateCarrollBuildings(db *gorm.DB) {
	simp := Building{Name: "Simperman", Latitude: 46.599809, Longitude: -112.037425}
	ocon := Building{Name: "O'Connell", Latitude: 46.600769, Longitude: -112.04027}
	stch := Building{Name: "St. Charles", Latitude: 46.600854, Longitude: -112.038613}
	chapel := Building{Name: "All Saints' Chapel", Latitude: 46.601338, Longitude: -112.038511}
	libr := Building{Name: "Library", Latitude: 46.601888, Longitude: -112.0381}
	stal := Building{Name: "St. Albert", Latitude: 46.600557, Longitude: -112.037722}
	borro := Building{Name: "Borromeo", Latitude: 46.599731, Longitude: -112.03917}
	cube := Building{Name: "Cube", Latitude: 46.599888, Longitude: -112.040354}
	ceng := Building{Name: "Civil Engineering", Latitude: 46.599567, Longitude: -112.036194}
	pccc := Building{Name: "Canine Center", Latitude: 46.599641, Longitude: -112.035523}

	db.Create(&simp)
	db.Create(&ocon)
	db.Create(&stch)
	db.Create(&chapel)
	db.Create(&libr)
	db.Create(&stal)
	db.Create(&borro)
	db.Create(&cube)
	db.Create(&ceng)
	db.Create(&pccc)
}
