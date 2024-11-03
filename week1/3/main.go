package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Define item prices
var itemsPrice = 8000000

func main() {
	// Split the string from input, by "," and store in a slice
	inputSlice := inputToList()

	// Call the function and get the counts
	shopList := countShopList(inputSlice)

	// Evaluate all grouping methods and find the one with the minimum total amount
	minCost := evaluateGroupingMethods(itemsPrice, shopList)

	// Print the lowest cost
	fmt.Println(minCost)
}

// Get string of shopping list separated by ","
func inputToList() []string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	textList := strings.Split(scanner.Text(), ",")
	return textList
}

// Calculate shopping list with item counts
func countShopList(input []string) map[string]int {
	shopList := make(map[string]int)
	for _, item := range input {
		shopList[item]++
	}
	return shopList
}

// This function evaluates multiple grouping methods and finds the one with the lowest total price
func evaluateGroupingMethods(itemPrice int, shopListItemCounts map[string]int) int {

	var methods []int
	// Try different grouping methods
	result := findCombinations(15)
	for _, combination := range result {
		methodCost := calculateGroups(itemPrice, shopListItemCounts, combination)
		methods = append(methods, methodCost)
	}
	return min(methods...)
}

// Function to find all combinations of numbers from 2 to 5 that sum up to 'target'
func findCombinations(target int) [][]int {
	// Base case: if the target is 0, return an empty combination
	if target == 0 {
		return [][]int{{}}
	}

	var result [][]int
	// Try numbers from 2 to 5
	for i := 2; i <= 5; i++ {
		if target >= i {
			// Get combinations for the remaining target
			subCombinations := findCombinations(target - i)
			// Add current number 'i' to each of the sub-combinations
			for _, comb := range subCombinations {
				result = append(result, append([]int{i}, comb...))
			}
		}
	}
	return result
}

// Function to calculate total price based on grouping
func calculateGroups(itemPrice int, shopListItemCounts map[string]int, groupSizes []int) int {
	// Create a copy of the original map
	shopListCopy := make(map[string]int)
	for item, count := range shopListItemCounts {
		shopListCopy[item] = count
	}

	totalCost := 0

	for _, groupSize := range groupSizes {
		group := make(map[string]int)

		// Add unique items to group (up to groupSize)
		for item, count := range shopListCopy {
			if len(group) < groupSize && count > 0 {
				group[item] = 1
				shopListCopy[item]--
			}
		}

		// Apply discount to the group and add to total cost
		totalCost += calculateDiscount(itemPrice, group)
	}

	// Add remaining items at full price
	for _, count := range shopListCopy {
		totalCost += count * itemPrice
	}

	return totalCost
}

// Function to calculate total price based on the number of unique items in a group
func calculateDiscount(itemPrice int, group map[string]int) int {
	totalPrice := len(group) * itemPrice

	// Apply discount based on the number of unique items
	switch len(group) {
	case 5:
		totalPrice = totalPrice - totalPrice*25/100
	case 4:
		totalPrice = totalPrice - totalPrice*20/100
	case 3:
		totalPrice = totalPrice - totalPrice*10/100
	case 2:
		totalPrice = totalPrice - totalPrice*5/100
	}

	return totalPrice
}

// Helper function to find the minimum value among multiple integers
func min(values ...int) int {
	minVal := values[0]
	for _, value := range values {
		if value < minVal {
			minVal = value
		}
	}
	return minVal
}
