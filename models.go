package main

import (
	"gorm.io/gorm"
)

// Room represents a physical room where classes are held.
type Room struct {
	gorm.Model
	Name       string   // Room name
	BuildingID uint     // Foreign key for the building (building ID)
	Building   Building // The building associated with this room
	Latitude   float64
	Longitude  float64
	Classes    []Class // Classes held in this room
}

// Class represents a university class.
type Class struct {
	gorm.Model
	Name      string
	Code      string `gorm:"unique"`
	RoomID    uint   // Foreign key for the room (room ID)
	Room      Room   // The room where this class is held
	Professor string
}

// Building represents a building where rooms are located.
type Building struct {
	gorm.Model
	Name      string
	Latitude  float64
	Longitude float64
	Rooms     []Room // Rooms in this building
}
