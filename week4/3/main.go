package main

import (
	"errors"
	"sync"
	"time"
)

// Thermostat represents the thermostat structure with current and target temperatures
type Thermostat struct {
	CurrentTemperature int
	TargetTemperature  int
}

// Room represents a room with an ID, thermostat, occupancy status, and fan status
type Room struct {
	ID         string
	Thermostat Thermostat
	Occupied   bool
	FanRunning bool
	// mu         sync.Mutex // For concurrent access
}

// SystemController manages a set of rooms
type SystemController struct {
	Rooms map[string]*Room
	mu    sync.Mutex
}

// NewSystemController creates a new SystemController instance
func NewSystemController() *SystemController {
	return &SystemController{
		Rooms: make(map[string]*Room),
	}
}

// AddRoom adds a new room to the system
func (s *SystemController) AddRoom(room *Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if room.ID == "" {
		return errors.New("room ID is required")
	}
	if _, exists := s.Rooms[room.ID]; exists {
		return errors.New("room already exists")
	}
	s.Rooms[room.ID] = room
	return nil
}

// UpdateRoomTemperature updates the temperature of a specific room
func (s *SystemController) UpdateRoomTemperature(roomID string, newTemp int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if newTemp < 0 {
		return errors.New("invalid target temperature")
	}
	room, exists := s.Rooms[roomID]
	if !exists {
		return errors.New("room does not exist")
	}

	room.Thermostat.CurrentTemperature = newTemp
	if room.Thermostat.CurrentTemperature == room.Thermostat.TargetTemperature {
		return room.StopFan()
	} else {
		return room.StartFan()
	}
}

// GenerateReports generates a report of all rooms and their fan status
func (s *SystemController) GenerateReports() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	report := make(map[string]string)
	for _, room := range s.Rooms {
		status := "off"
		if room.FanRunning {
			if room.Thermostat.CurrentTemperature > room.Thermostat.TargetTemperature {
				status = "cooling"
			} else if room.Thermostat.CurrentTemperature < room.Thermostat.TargetTemperature {
				status = "heating"
			}
		}
		report[room.ID] = status
	}
	return report
}

// SetOccupancy sets the occupancy status of the room
func (r *Room) SetOccupancy(occupied bool) {

	r.Occupied = occupied
	if !r.Occupied {
		_ = r.StopFan()
	}
	if r.Occupied {
		_ = r.StartFan()
	}
}

// StartFan starts the fan in the room
func (r *Room) StartFan() error {
	if !r.Occupied {
		return errors.New("room is not occupied")
	}
	if r.FanRunning {
		return errors.New("fan already running")
	}
	if r.Thermostat.CurrentTemperature == r.Thermostat.TargetTemperature {
		return errors.New("no adjustment needed")
	}

	r.FanRunning = true

	go func() {
		for r.FanRunning && r.Thermostat.CurrentTemperature != r.Thermostat.TargetTemperature {
			if r.Thermostat.CurrentTemperature > r.Thermostat.TargetTemperature {
				r.Thermostat.CurrentTemperature--
			} else if r.Thermostat.CurrentTemperature < r.Thermostat.TargetTemperature {
				r.Thermostat.CurrentTemperature++
			}
			// Simulate 1°C change per second
			time.Sleep(1 * time.Second)
		}

		// Stop the fan once the target temperature is reached
		r.StopFan()
	}()

	return nil
}

// StopFan stops the fan in the room
func (r *Room) StopFan() error {

	if !r.FanRunning {
		return errors.New("fan not running")
	}
	r.FanRunning = false
	return nil
}

// GetCurrentTemperature returns the current temperature of the room
func (r *Room) GetCurrentTemperature() int {

	return r.Thermostat.CurrentTemperature
}

// GetTargetTemperature returns the target temperature of the room
func (r *Room) GetTargetTemperature() int {

	return r.Thermostat.TargetTemperature
}

// GetIsRoomOccupied returns the occupancy status of the room
func (r *Room) GetIsRoomOccupied() bool {

	return r.Occupied
}

// GetIsFanRunning returns whether the fan is running
func (r *Room) GetIsFanRunning() bool {

	return r.FanRunning
}

// func main() {
// 	controller := NewSystemController()

// 	// Example: Adding rooms and updating temperatures
// 	room1 := &Room{
// 		ID: "101",
// 		Thermostat: Thermostat{
// 			CurrentTemperature: 25,
// 			TargetTemperature:  22,
// 		},
// 		Occupied:   true,
// 		FanRunning: false,
// 	}

// 	room2 := &Room{
// 		ID: "102",
// 		Thermostat: Thermostat{
// 			CurrentTemperature: 20,
// 			TargetTemperature:  22,
// 		},
// 		Occupied:   true,
// 		FanRunning: false,
// 	}

// 	controller.AddRoom(room1)
// 	controller.AddRoom(room2)

// 	controller.UpdateRoomTemperature("101", 23)
// 	controller.UpdateRoomTemperature("102", 19)

// 	report := controller.GenerateReports()
// 	for roomID, status := range report {
// 		fmt.Printf("%s -> %s\n", roomID, status)
// 	}

// }
