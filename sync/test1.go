package main

import "fmt"

func main() {
	ch := make(chan int)

	for i := 1; i <= 10; i++ {
		// i := i
		go func() {
			ch <- i * i
		}()
	}

	sum := 0
	for i := 0; i < 10; i++ {
		v := <-ch
		sum += v
	}

	fmt.Println(sum)
}
