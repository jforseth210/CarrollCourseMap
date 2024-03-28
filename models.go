package main

import (
	"time"

	"gorm.io/gorm"
)

// Room represents a physical room where classes are held.
type Room struct {
	gorm.Model
	Name       string   // Room name
	BuildingID uint     // Foreign key for the building (building ID)
	Building   Building // Foreign key for the building (building ID)
	Latitude   float64
	Longitude  float64
}

// Professor represents a university professor.
type Professor struct {
	gorm.Model
	Name    string   // Professor's name
	Classes []*Class `gorm:"many2many:professor_classes;"` // Many-to-many relationship with classes
}

// Class represents a university class.
type Class struct {
	gorm.Model
	Name         string       // Class name
	Description  string       // Class description
	RoomID       uint         // Foreign key for the room (room ID)
	Room         Room         // Room associated with the class
	Professors   []*Professor `gorm:"many2many:professor_classes;"` // Many-to-many relationship with professors
	ClassPeriods []ClassPeriod
}

type ClassPeriod struct {
	gorm.Model
	ClassID   uint
	Class     Class
	StartTime time.Time
	EndTime   time.Time
}

// Building represents a building where rooms are located.
type Building struct {
	gorm.Model
	Name      string
	Latitude  float64
	Longitude float64
}
