package main

import (
	"fmt"
	"love-signal-users/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("config: %v", cfg)
}
