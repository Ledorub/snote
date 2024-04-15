package main

import (
	"fmt"
	"github.com/ledorub/snote-api/internal/config"
	"github.com/ledorub/snote-api/internal/logger"
)

func main() {
	cfg := config.New()
	fmt.Printf("Got port %d.\n", cfg.Port)
	log := logger.New()
	log.Println("Set up logger.")
}
