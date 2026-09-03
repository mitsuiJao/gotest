package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()

	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)

		for i := range in {
			out <- i * i
		}
	}()

	return out
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))

	for _, c := range cs {
		go func() {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	arr := []int{}
	for i := 1; i <= 30; i++ {
		arr = append(arr, i)
	}

	in := gen(arr...)

	outs := make([]<-chan int, 3)
	for i := range 3 {
		outs[i] = square(in)
	}

	var result int = 0
	out := merge(outs...)
	for v := range out {
		n := v
		result += n
	}

	fmt.Println(result)
}
