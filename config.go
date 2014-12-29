package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	// Network timeout in seconds
	Timeout int
	// SMTP relay config
	SMTP SMTPConfig
	// Alert properties
	Alert   Alert
	Targets []Target
}

type Alert struct {
	// On alert, send to this email address
	ToEmail string
	// On alert, send from this email address
	FromEmail string
	// Trigger an alert every x seconds when in failed state
	Interval int
}

type SMTPConfig struct {
	Hostname string
	Port     int
}

// Opening (or creating) config file in JSON format
func readConfig(filename string) Config {
	config := Config{
		Timeout: 10,
		Targets: []Target{Target{Name: "Local HTTP Server", Addr: "http://localhost"}},
	}

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		// unaccessible or not exisiting file -> creatoin
		file, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}

		// config file just created
		err := json.NewEncoder(file).Encode(config)
		if err != nil {
			log.Fatal(err)
		}

	} else {

		err = json.NewDecoder(file).Decode(&config)

		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}
