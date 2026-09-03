package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func merge(ctx context.Context, a, b <-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-a:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-b:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
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

	for i := 1; i <= 20; i++ {
		a <- i
		time.Sleep(200 * time.Millisecond)
	}
}
func genB(b chan<- int) {
	defer close(b)

	for i := 100; i <= 119; i++ {
		b <- i
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	a := make(chan int)
	b := make(chan int)

	go genA(a)
	go genB(b)

	merged := merge(ctx, a, b)

	var count int
	for v := range merged {
		_ = v
		count++
	}

	fmt.Println(count)
}
