package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Struct to hold location data
type Location struct {
	Name       string
	X, Y       int
	From       string
	XOffset    int
	YOffset    int
	Calculated bool
}

func main() {
	input := readInput()
	locationMap, order := parseInput(input)
	printLocations(locationMap, order)
}

// Read input until an empty line
func readInput() []string {
	scanner := bufio.NewScanner(os.Stdin)
	var input []string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		input = append(input, line)
	}
	return input
}

// Parse input to determine locations
func parseInput(input []string) (map[string]Location, []string) {
	locationMap := make(map[string]Location)
	var order []string // To maintain the order of locations

	// First pass: Register start location and other locations
	for _, line := range input {
		if strings.HasPrefix(line, "start") {
			var startX, startY int
			fmt.Sscanf(line, "start x=%d y=%d", &startX, &startY)
			locationMap["start"] = Location{Name: "start", X: startX, Y: startY, Calculated: true} // Mark start as calculated
		} else {
			var name, from string
			var x, y int
			fmt.Sscanf(line, "%s from %s x=%d y=%d", &name, &from, &x, &y)

			// for prevent duplication
			if exists(order, name) {
				continue
			}

			// Store the location with its offsets and dependency
			locationMap[name] = Location{Name: name, From: from, XOffset: x, YOffset: y, Calculated: false}
			order = append(order, name) // Record the order of the location
		}
	}

	// Second pass: Recursively calculate positions for each location
	for _, name := range order {
		resolveLocation(locationMap, name)
	}

	return locationMap, order
}

// Resolve the position of a location recursively
func resolveLocation(locationMap map[string]Location, name string) Location {
	loc := locationMap[name]

	// If already calculated, return the location
	if loc.Calculated {
		return loc
	}

	// If the location depends on another location, resolve the dependency first
	if fromLoc, exists := locationMap[loc.From]; exists {
		// Recursively resolve the location it's dependent on
		fromLoc = resolveLocation(locationMap, loc.From)

		// Calculate this location's absolute position
		loc.X = loc.XOffset + fromLoc.X
		loc.Y = loc.YOffset + fromLoc.Y
		loc.Calculated = true
		locationMap[name] = loc // Update the map with calculated values
	}

	return loc
}

// Print the locations in the order they were added
func printLocations(locationMap map[string]Location, order []string) {
	for _, name := range order {
		if name != "start" { // Skip printing "start"
			loc := locationMap[name]
			fmt.Printf("%s x=%d y=%d\n", loc.Name, loc.X, loc.Y)
		}
	}
}

// Function to check if a name exists in a list
func exists(list []string, name string) bool {
	for _, item := range list {
		if item == name {
			return true
		}
	}
	return false
}
