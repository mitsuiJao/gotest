package main

import "fmt"

func worker(w int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		results <- job * 3
	}
}

func main() {
	jobs := make(chan int, 20)
	results := make(chan int, 20)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for i := 1; i <= 20; i++ {
		jobs <- i
	}
	close(jobs)

	sum := 0
	for j := 1; j <= 20; j++ {
		v := <-results
		sum += v
	}

	fmt.Println(sum)
}
