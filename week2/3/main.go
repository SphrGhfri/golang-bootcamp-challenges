package main

import (
	"errors"
	"sync"
)

type Airpod struct {
	state   string // Combined state: "Docked", "Connected", or "Disconnected"
	channel chan byte
}

type AirpodCase struct {
	left        *Airpod
	right       *Airpod
	isConnected bool
	channel     chan byte
	mu          sync.Mutex
}

func NewAirpodCase() *AirpodCase {
	return &AirpodCase{
		left: &Airpod{
			state:   "Docked", // Default state is "Docked"
			channel: make(chan byte, 500),
		},
		right: &Airpod{
			state:   "Docked", // Default state is "Docked"
			channel: make(chan byte, 500),
		},
		isConnected: false,
		channel:     make(chan byte, 500),
	}
}

// GetRightAirpod returns the right AirPod
func (a *AirpodCase) GetRightAirpod() *Airpod {
	return a.right
}

// GetLeftAirpod returns the left AirPod
func (a *AirpodCase) GetLeftAirpod() *Airpod {
	return a.left
}

// GetState returns the current state of the AirPod
func (a *Airpod) GetState() string {
	return a.state
}

// UndockLeft undocks the left AirPod and adjusts the connection state
func (a *AirpodCase) UndockLeft() *Airpod {
	a.mu.Lock()
	defer a.mu.Unlock()
	leftPod := a.GetLeftAirpod()

	if leftPod.GetState() == "Docked" {
		leftPod.state = "Disconnected" // Initially set to Disconnected
		if a.isConnected {
			leftPod.state = "Connected"
			go func() {
				for msg := range a.channel {
					a.mu.Lock()
					leftPod.channel <- msg

					a.mu.Unlock()
				}
			}()
		}
		return leftPod
	}
	return nil
}

// UndockRight undocks the right AirPod and adjusts the connection state
func (a *AirpodCase) UndockRight() *Airpod {
	a.mu.Lock()
	defer a.mu.Unlock()
	rightPod := a.GetRightAirpod()

	if rightPod.GetState() == "Docked" {
		rightPod.state = "Disconnected" // Initially set to Disconnected
		if a.isConnected {
			rightPod.state = "Connected"
			go func() {
				for msg := range a.channel {
					a.mu.Lock()
					rightPod.channel <- msg

					a.mu.Unlock()
				}
			}()
		}
		return rightPod
	}
	return nil
}

// DockLeft docks the left AirPod
func (a *AirpodCase) DockLeft() error {
	leftPod := a.GetLeftAirpod()
	if leftPod.GetState() == "Docked" {
		return errors.New("left AirPod is already docked")
	}

	leftPod.state = "Docked"
	leftPod.channel = make(chan byte, 500) // Reset the channel when docked
	return nil
}

// DockRight docks the right AirPod
func (a *AirpodCase) DockRight() error {
	rightPod := a.GetRightAirpod()

	if rightPod.GetState() == "Docked" {
		return errors.New("right AirPod is already docked")
	}

	rightPod.state = "Docked"
	rightPod.channel = nil // Reset the channel when docked
	return nil
}

// GetChannel returns the current channel only if the AirPod is connected
func (a *Airpod) GetChannel() chan byte {

	if a.state == "Connected" {
		return a.channel
	}
	return nil
}

// ConnectBluetooth connects both AirPods and the case in a specific order
func (c *AirpodCase) ConnectBluetooth(ch chan byte) error {

	if c.isConnected {
		return errors.New("already connected")
	}

	c.isConnected = true

	if c.left.state == "Disconnected" {
		c.left.state = "Connected"
	}
	if c.right.state == "Disconnected" {
		c.right.state = "Connected"
	}

	go func() {
		for msg := range ch {
			// c.channel <- msg
			c.mu.Lock()
			if c.left.state == "Connected" {
				c.left.channel <- msg
			}
			if c.right.state == "Connected" {
				c.right.channel <- msg
			}
			c.mu.Unlock()
		}
	}()

	return nil
}
