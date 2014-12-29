package main

import (
	"log"
	"net"
	"time"
)

type Target struct {
	// Name of the Target
	Name string
	// Address (ex: "localhost:80" of the target
	Addr string
	// Polling interval, in seconds
	Interval int
}

type TargetStatus struct {
	Target    *Target
	Online    bool
	Since     time.Time
	LastCheck time.Time
}

func startTarget(t Target, res chan TargetStatus) {
	go runTarget(t, res)
}

func runTarget(t Target, res chan TargetStatus) {

	log.Println("starting runtarget on ", t.Name)
	if t.Interval == 0 {
		t.Interval = 1
	}
	ticker := time.Tick(time.Duration(t.Interval) * time.Second)
	for {
		status := TargetStatus{Target: &t}
		// Polling
		conn, err := net.Dial("tcp", t.Addr)

		if err != nil {
			if status.Online {
				status.Since = time.Now()
			}
			// Error during connect
			status.Online = false
		} else {
			if !status.Online {
				status.Since = time.Now()
			}
			// Connect ok
			status.Online = true
			conn.Close()
		}
		status.LastCheck = time.Now()

		res <- status

		// waiting for ticker
		<-ticker
	}
}
