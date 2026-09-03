package main

import (
	"fmt"
	"sync"
)

func worker(jobs <-chan int, counts map[string]int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		mu.Lock()
		if job%2 == 0 {
			counts["even"]++
		} else {
			counts["odd"]++
		}
		mu.Unlock()
	}
}

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	jobs := make(chan int)

	counts := make(map[string]int)

	for w := 0; w < 5; w++ {
		wg.Add(1)
		go worker(jobs, counts, &mu, &wg)
	}

	for i := 0; i < 100; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	fmt.Printf("odd: %d, even: %d\n", counts["odd"], counts["even"])
}
