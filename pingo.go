package main

import (
	"flag"
	"log"
)

// Init config

// Main function
func main() {
	filename := flag.String("f", "config.json", "JSON configuration file")
	httpPort := flag.Int("p", 8888, "HTTP port")

	flag.Parse()

	// Config
	log.Printf("Opening config file: %s\n", *filename)
	config := readConfig(*filename)
	log.Println("Config loaded")

	// Running
	res := make(chan TargetStatus)
	state := NewState()

	for _, target := range config.Targets {
		if target.Addr != "" {
			startTarget(target, res, config)
		}
	}

	// HTTP
	go startHttp(*httpPort, state)

	for {
		select {
		case status := <-res:
			state.Lock()
			state.State[status.Target] = status
			state.Unlock()
		}
	}
}
