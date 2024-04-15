package main

import (
	"fmt"
	"github.com/ledorub/snote-api/internal/config"
)

func main() {
	cfg := config.New()
	fmt.Printf("Got port %d.", cfg.Port)
}
