package main

import (
	"fmt"
	"sync"
	"time"
)

func merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for v := range a {
			out <- v
		}
	}()

	go func() {
		defer wg.Done()
		for v := range b {
			out <- v
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func genA(a chan<- int) {
	defer close(a)

	for i := 1; i <= 5; i++ {
		a <- i
		time.Sleep(100 * time.Millisecond)
	}
}
func genB(b chan<- int) {
	defer close(b)

	for i := 100; i <= 104; i++ {
		b <- i
		time.Sleep(150 * time.Millisecond)
	}
}

func main() {
	a := make(chan int)
	b := make(chan int)

	go genA(a)
	go genB(b)

	merged := merge(a, b)

	for v := range merged {
		fmt.Println(v)
	}

}
