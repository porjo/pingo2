// The MIT License (MIT)
//
// Copyright (c) 2014 Ian Bishop
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Pingo2 is an open source tool written in Golang that allows you to monitor the avilability of TCP server applications.
package main

import (
	"flag"
	"log"
)

var debug = false

// Init config

// Main function
func main() {
	//filename := flag.String("f", "config.toml", "TOML configuration file")
	filename := flag.String("f", "config.json", "JSON configuration file")
	httpPort := flag.Int("p", 8888, "HTTP port")
	flag.BoolVar(&debug, "d", false, "Enable debug output")

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
