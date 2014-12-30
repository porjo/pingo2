package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
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

// Opening (or creating) config file in TOML format
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
		err := toml.NewEncoder(file).Encode(config)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		_, err := toml.DecodeReader(file, &config)

		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}
