package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// Init config

var filename = flag.String("f", "config.json", "JSON configuration file")
var httpPort = flag.Int("p", 8888, "HTTP port")

func startBrowser(port int, url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows", "darwin":
		err = exec.Command("cmd", "/c", "start", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// Main function
func main() {

	flag.Parse()

	// Config
	log.Println("Opening config file: ", *filename)
	config := readConfig(*filename)
	log.Printf("Config loaded")

	// Running
	res := make(chan TargetStatus)
	end := make(chan int)
	state := NewState()

	for _, target := range config.Targets {
		startTarget(target, res, end)
	}

	// HTTP	

	go startHttp(*httpPort, state)
	go startBrowser(*httpPort, fmt.Sprintf("http://localhost:%d/status", *httpPort))

	for {
		select {
		case <-end:
			log.Println("One of the checker ended...")
		case status := <-res:
			state.Lock.Lock()
			if s, ok := state.State[status.Target]; ok {
				if s.Online != status.Online {
					s.Online = status.Online
					s.Since = status.Since
				}
				s.LastCheck = status.Since
				status = s
			} else {
				status.LastCheck = status.Since
			}
			state.State[status.Target] = status

			state.Lock.Unlock()
		}
	}

}
