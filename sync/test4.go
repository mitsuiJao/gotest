package main

import (
	"fmt"
	"time"
)

func main() {
	chA := make(chan int)
	chB := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			chA <- i
			time.Sleep(300 * time.Millisecond)
		}
	}()

	go func() {
		for j := 100; j <= 104; j++ {
			chB <- j
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for n := 0; n < 10; n++ {
		select {
		case a := <-chA:
			fmt.Printf("A: %d\n", a)
		case b := <-chB:
			fmt.Printf("B: %d\n", b)
		}
	}
}
