package main

import (
	"fmt"
)

var (
	previousX = 1
	previousY = 1
)

func FindTreasure(ch chan interface{}) {
	// Step 1: Read initial values from the channel
	maxQuery := (<-ch).(int)
	xMapMax := (<-ch).(int) // number of rows
	yMapMax := (<-ch).(int) // number of columns
	xMapMin := 1
	yMapMin := 1
	x0 := 1
	y0 := 1
	direction := 1
	// answer := "not found yet"
	// x_answer := 0
	// y_answer := 0

	xMax := xMapMax
	yMax := yMapMax
	xMin := xMapMin - 1
	yMin := yMapMin - 1

	var x_answer int
	// var y_answer int
	// }
	// func main() {

	// Send the mid-point coordinates
	ch <- x0
	ch <- y0

	// Receive the feedback from the detector
	feedback := (<-ch).(string)

	fmt.Println(feedback)

	// Searching along the x-axis
	for i := 0; i < maxQuery; i++ {
		x0 += direction * int(float64(xMax-xMin)/2)
		ch <- x0
		ch <- y0
		result := <-ch

		if result != nil {
		}

		// fmt.Println("resultX", result)

		x1 := x0 + 1
		ch <- x1
		ch <- y0
		result1 := <-ch

		x2 := x1 - 2
		ch <- x2
		ch <- y0
		result2 := <-ch

		// if result == "treasure found" {
		// 	x_answer := x0
		// 	y_answer := y0
		// 	fmt.Printf("answer = %d,%d\n", x_answer, y_answer)
		// 	break
		// } else if result1 == "treasure found" {
		// 	x_answer := x1
		// 	y_answer := y0
		// 	fmt.Printf("answer = %d,%d\n", x_answer, y_answer)
		// 	break
		// } else if result2 == "treasure found" {
		// 	x_answer := x2
		// 	y_answer := y0
		// 	fmt.Printf("answer = %d,%d\n", x_answer, y_answer)
		// 	break
		// } else
		if result1 == "$s$s$" {
			direction = 1
			xMin = x0
		} else if result2 == "$s$s$" {
			direction = -1
			xMax = x0
		} else if result1 == "s.s.s" && result2 == "same as before" {
			x_answer = x0
			fmt.Printf("answer found: location = %d\n", x_answer)
			break
		}

		if x1 == xMapMax && result2 == "s.s.s" {
			x_answer = x1
			fmt.Printf("answer found: location = %d\n", x_answer)
			break
		}

		previousX = x0
		// fmt.Printf("run %d:\n direction = %d,\n x_min = %d,\n x_max = %d,\n x0 = %d, \n x1 = %d,\n x2 = %d,\n", i, direction, xMin, xMax, x0, x1, x2)
	}

	fmt.Println("________________")

	direction = 1

	// Searching along the y-axis
	for i := 0; i < yMapMax; i++ {
		y0 += direction * int(float64(yMax-yMin)/2)

		ch <- x_answer
		ch <- y0
		result := <-ch
		if result != nil {
		}
		// fmt.Println("resultX", result)

		y1 := y0 + 1
		ch <- x_answer
		ch <- y1
		result1 := <-ch

		y2 := y1 - 2
		ch <- x_answer
		ch <- y2
		result2 := <-ch

		// if result == "treasure found" {
		// 	y_answer = y0
		// 	fmt.Printf("answer = %d,%d\n", x_answer, y_answer)
		// 	break
		// } else if result1 == "treasure found" {
		// 	y_answer = y1
		// 	fmt.Printf("answer = %d,%d\n", x_answer, y_answer)
		// 	break
		// } else if result2 == "treasure found" {
		// 	y_answer = y2
		// 	fmt.Printf("answer = %d,%d\n", x_answer, y_answer)
		// 	break
		// } else

		if result1 == "$s$s$" {
			direction = 1
			yMin = y0
		} else if result2 == "$s$s$" {
			direction = -1
			yMax = y0
		}
		// else if result1 == "s.s.s" && result2 == "same as before" {
		// 	y_answer := y0
		// 	// fmt.Printf("answer found: location = %d,%d\n", x_answer, y_answer)
		// 	break
		// }

		// if y1 == yMapMax && result2 == "s.s.s" {
		// 	y_answer := y1
		// 	// fmt.Printf("answer = %d,%d\n", x_answer, y_answer)
		// 	break
		// }

		previousY = y0
		// fmt.Printf("run %d:\n direction = %d,\n y_min = %d,\n y_max = %d,\n y0 = %d, \n y1 = %d,\n y2 = %d,\n", i, direction, yMin, yMax, y0, y1, y2)
	}
}
