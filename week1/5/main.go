package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Define the time format
var timeLayout = "15:04" // 24-hour format

var restaurantAverageMap = make(map[string]int)

func main() {
	restaurantsInput, ordersInput, request := readInput() // Read input lines
	// Parse all restaurants average time
	for _, input := range restaurantsInput {
		restaurant, average := parsRestaurant(input)
		restaurantAverageMap[restaurant] = average
	}

	if request == "overlap" {
		parsedOrders := makeOrdersList(ordersInput, false)
		// Create a map to store the overlap counts for each minute of the day
		minuteCounts := make(map[string]int)

		// Loop through each minute from 00:00 to 23:59
		for d := 1; d < 3; d++ {
			for h := 0; h < 24; h++ {
				for m := 0; m < 60; m += 1 { // increment by 1 to get minutes
					checkTimer := time.Date(0, 1, d, h, m, 0, 0, time.UTC) // Specify the date as needed

					// Count how many orders overlap with the current minute
					for _, order := range parsedOrders {
						if isOverlap(checkTimer, order.received, order.delivered) {
							minuteCounts[checkTimer.Format("15:04")]++ // Format as HH:MM for the map key
						}
					}
				}
			}
		}
		// Find the maximum overlap count
		maxCount := 0
		for _, count := range minuteCounts {
			if count > 1 && count > maxCount {
				maxCount = count
			}
		}
		// fmt.Println(minuteCounts)
		// Output the maximum number of overlaps
		fmt.Println(maxCount)
	} else {
		// Parse all orders
		parsedOrders := makeOrdersList(ordersInput, true)
		delayCount := 0
		for _, order := range parsedOrders {
			restaurant := order.restaurant
			received := order.received
			delivered := order.delivered
			if restaurant == request {
				restaurantAverageTime := restaurantAverageMap[request]
				duration := delivered.Sub(received)
				if (int(duration.Minutes()) - 15) >= restaurantAverageTime {
					delayCount++
				}
			}
		}
		fmt.Println(delayCount)
	}

}
func readInput() ([]string, []string, string) {
	scanner := bufio.NewScanner(os.Stdin)
	var restaurantsInput []string
	var ordersInput []string

	scanner.Scan()                           // Read the first line to get the count of Restaurants and their average time
	count, _ := strconv.Atoi(scanner.Text()) // n number of Restaurants

	// Read the specified number of input lines
	for i := 0; i < count; i++ {
		if scanner.Scan() {
			restaurantsInput = append(restaurantsInput, strings.TrimSpace(scanner.Text())) // Append each line
		}
	}
	scanner.Scan()
	count, _ = strconv.Atoi(scanner.Text()) // m number of Restaurants

	// Read the specified number of input lines
	for i := 0; i < count; i++ {
		if scanner.Scan() {
			ordersInput = append(ordersInput, strings.TrimSpace(scanner.Text())) // Append each line
		}
	}
	scanner.Scan()
	requestType := strings.TrimSpace(scanner.Text())

	return restaurantsInput, ordersInput, requestType
}

func detectCustomerRestaurant(a string, b string) (string, string, bool) {
	// Regular expression to match from the start until the first number
	_, exists1 := restaurantAverageMap[a]
	_, exists2 := restaurantAverageMap[b]
	if exists1 {
		return b, a, false
	} else if exists2 {
		return a, b, true
	} else {
		fmt.Println("NOT FOUND!")
		return "", "", false
	}
}

func parsRestaurant(input string) (restaurantName string, averageTime int) {
	// Regular expression to capture the string part and integer part
	re := regexp.MustCompile(`^(.*?)(\d+)$`)
	matches := re.FindStringSubmatch(input)
	restaurantName = strings.TrimSpace(matches[1])
	averageTime, _ = strconv.Atoi(matches[2])
	return restaurantName, averageTime
}

func parseOrder(input string, forDelay bool) (restaurantName string, receivedTime time.Time, deliveredTime time.Time, customerName string) {
	// Regular expression to capture the restaurant name, received time, delivered time, and customer name
	re := regexp.MustCompile(`^(.*?)\s+(\d{1,2}:\d{1,2})\s+(\d{1,2}:\d{1,2})\s+(.+)$`)
	matches := re.FindStringSubmatch(input)
	customerName, restaurantName, isFirstName := detectCustomerRestaurant(strings.TrimSpace(matches[1]), strings.TrimSpace(matches[4]))
	if isFirstName {
		deliveredTime, _ = time.Parse(timeLayout, matches[2])
		receivedTime, _ = time.Parse(timeLayout, matches[3])
		if forDelay {

			newTime := receivedTime.Add(-5 * time.Minute)
			if newTime.Year() < 0 {
				receivedTime = receivedTime.Add(24 * time.Hour)
				receivedTime = receivedTime.Add(-5 * time.Minute)
			} else {
				receivedTime = receivedTime.Add(-5 * time.Minute)
			}
		} else {
			receivedTime, _ = time.Parse(timeLayout, matches[3])
		}
	} else if !isFirstName {
		// First part is restaurant's name
		receivedTime, _ = time.Parse(timeLayout, matches[2])
		deliveredTime, _ = time.Parse(timeLayout, matches[3])
	}
	if deliveredTime.Before(receivedTime) {
		deliveredTime = deliveredTime.Add(24 * time.Hour)
	}
	return restaurantName, receivedTime, deliveredTime, customerName
}

func makeOrdersList(ordersInput []string, forDelay bool) (output []struct {
	restaurant string
	customer   string
	received   time.Time
	delivered  time.Time
}) {
	var parsedOrders []struct {
		restaurant, customer string
		received, delivered  time.Time
	}
	// Parse all orders
	for i := 0; i < len(ordersInput); i++ {
		// Parsing order
		restaurant, received, delivered, customer := parseOrder(ordersInput[i], forDelay)

		// Appending parsed values to parsedOrders slice
		parsedOrders = append(parsedOrders, struct {
			restaurant, customer string
			received, delivered  time.Time
		}{
			restaurant: restaurant,
			customer:   customer,
			received:   received,
			delivered:  delivered,
		})
	}
	return parsedOrders
}

func isOverlap(checkTime, received, delivered time.Time) bool {
	return (checkTime.After(received) || checkTime.Equal(received)) && (checkTime.Before(delivered) || checkTime.Equal(received))
}
