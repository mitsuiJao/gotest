package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, canacel := context.WithTimeout(context.Background(), 3*time.Second)
	defer canacel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("cancelled:", ctx.Err())
				return
			case <-time.After(3 * time.Second):
				fmt.Println("stop:")
			}
		}
	}()

	time.Sleep(5 * time.Second)
}
