package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func trap(height []int) int {
	if len(height) == 0 {
		return 0
	}

	left, right := 0, len(height)-1
	leftMax, rightMax := height[left], height[right]
	waterTrapped := 0

	for left < right {
		if leftMax < rightMax {
			left++
			if height[left] < leftMax {
				waterTrapped += leftMax - height[left]
			} else {
				leftMax = height[left]
			}
		} else {
			right--
			if height[right] < rightMax {
				waterTrapped += rightMax - height[right]
			} else {
				rightMax = height[right]
			}
		}
	}

	return waterTrapped
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	nStr, _ := reader.ReadString('\n')
	n, _ := strconv.Atoi(strings.TrimSpace(nStr))

	heightsInput, _ := reader.ReadString('\n')

	heightStrs := strings.Split(strings.TrimSpace(heightsInput), ",")
	height := make([]int, n)

	for i, h := range heightStrs {
		height[i], _ = strconv.Atoi(strings.TrimSpace(h)) // Trim spaces around heights
	}

	result := trap(height)
	fmt.Println(result)
}
