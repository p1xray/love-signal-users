package main

import (
	"fmt"

	"github.com/p1xray/love-signal-users/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("config: %v", cfg)
}
