package main

import "fmt"

type Point struct {
	X int
	Y int
}

func main() {
	arr := []int{1, 2, 3, 4, 5}
	fmt.Println(arr)

	p := Point{X: 1, Y: 2}
	fmt.Println(p)
}
