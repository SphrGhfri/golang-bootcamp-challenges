package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	inputs := readInput() // Read input lines

	for _, input := range inputs {
		spaceCounter := 0
		result := ""                             // String to accumulate the result
		wordLen := utf8.RuneCountInString(input) // Length of the current input

		for _, char := range input {
			switch {
			case isNumber(char):
				// Convert number to ASCII value, scale by word length, and reverse the string
				numValue := int(char) * wordLen
				result += reverseString(fmt.Sprintf("%d", numValue))
				spaceCounter++
			case isSpace(char):
				result += strings.Repeat("_", spaceCounter) // Append underscores for spaces
				spaceCounter = 0                            // Reset counter
			default:
				result += string(encryptStringOrEmoji(char)) // Encrypt character and add to result
				spaceCounter++
			}
		}

		// Output the final result for the current input
		fmt.Println(reverseString(result))
	}
}

func readInput() []string {
	scanner := bufio.NewScanner(os.Stdin)
	var input []string

	scanner.Scan()                           // Read the first line to get the count of subsequent inputs
	count, _ := strconv.Atoi(scanner.Text()) // Convert to integer

	// Read the specified number of input lines
	for i := 0; i < count; i++ {
		if scanner.Scan() {
			input = append(input, scanner.Text()) // Append each line
		}
	}
	return input
}

// encryptStringOrEmoji encrypts a single character (rune)
func encryptStringOrEmoji(char rune) rune {
	if char >= 'a' && char <= 'z' {
		return 'a' + (char-'a'+(char%5))%26
	} else if char >= 'A' && char <= 'Z' {
		return 'A' + (char-'A'+(char%5))%26
	}
	return char + (char % 5) // Shift other characters
}

// isNumber checks if the character is a digit
func isNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

// isSpace checks if the character is a space
func isSpace(char rune) bool {
	return char == ' '
}

// reverseString reverses the given string
func reverseString(s string) string {
	runes := []rune(s) // Convert string to runes for Unicode support
	length := len(runes)

	// Reverse the runes in place
	for i := 0; i < length/2; i++ {
		runes[i], runes[length-1-i] = runes[length-1-i], runes[i]
	}

	return string(runes) // Convert runes back to a string
}
