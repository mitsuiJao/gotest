package main

import (
	"fmt"
	"sync"
)

type Config struct {
	APIKey string
}

var (
	config *Config
	once   sync.Once
)

func loadConfig() *Config {
	once.Do(func() {
		fmt.Println("config reading...")
		config = &Config{APIKey: "secret-key-12345"}
	})
	return config
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := loadConfig()
			fmt.Printf("goroutine %d: %+v\n", id, cfg)
		}(i)
	}

	wg.Wait()
}
