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
	state := NewState()

	for _, target := range config.Targets {
		startTarget(target, res)
	}

	// HTTP
	go startHttp(*httpPort, state)
	go startBrowser(*httpPort, fmt.Sprintf("http://localhost:%d/status", *httpPort))

	for {
		select {
		case status := <-res:
			state.Lock()
			state.State[status.Target] = status
			state.Unlock()
		}
	}

}
