## Pingo
=========

Pingo is an open source tool written in Golang that allows you to monitor the avilability of TCP server applications from a simple web view.

## Installation
================

To install:

	go get github.com/orcheus/pingo
	
Then build :
	go build github.com/orcheus/pingo
	
## Quick start
===========

Once Pingo has been built into an executable file, run it with no parameters.

By default, Pingo use the local "config.json" configuration file.
You can set a different configuration file using the parameter "-f" while starting Pingo.

If the configuration file doesn't existe, Pingo will create a default one, set to monitor the local port 80 every 5 seconds.

## Configuration
=============

Pingo can be configured using a json file.
This configuration file describes the different targets Pingo has to monitor.
For each target you can set a name, address (ex: "172.16.25.1:8080") and polling interval (in seconds).