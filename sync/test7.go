package main

func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func calculate(nums <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range nums {
			out <- n * 2
		}
	}()
	return out
}

func filter(chan<- int) <-chan int {

}

func main() {

}
