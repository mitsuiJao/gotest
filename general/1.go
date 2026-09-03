package main

import "fmt"

func main() {
	x := make(chan int)

	fmt.Printf("%T", x)
}
