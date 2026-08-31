package main

import "fmt"

func main() {
	ch := make(chan int, 5)

	for i := 1; i <= 5; i++ {
		go func() {
			ch <- i * 2
		}()
	}

	sum := 0
	for j := 1; j <= 5; j++ {
		v := <-ch
		sum += v
	}

	fmt.Println(sum)
}
