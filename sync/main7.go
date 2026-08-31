package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("cancelled:", ctx.Err())
				return
			default:
				time.Sleep(500 * time.Millisecond)
				fmt.Println("working...")
			}
		}
	}()

	time.Sleep(3 * time.Second)
}
