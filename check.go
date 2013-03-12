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
	Interval int64
}

type TargetStatus struct {
	Target    *Target
	Online    bool
	Since     time.Time
	LastCheck time.Time
}

func startTarget(t Target, res chan TargetStatus, end chan int) {
	go runTarget(t, res, end)
}

func runTarget(t Target, res chan TargetStatus, end chan int) {

	log.Println("starting runtarget on ", t.Name)
	if t.Interval == 0 {
		t.Interval = 1
	}
	ticker := time.Tick(time.Duration(t.Interval) * time.Second)
	for {
		// Polling

		var status TargetStatus
		conn, err := net.Dial("tcp", t.Addr)

		if err != nil {
			// Connect ok
			status = TargetStatus{Target: &t, Online: false, Since: time.Now()}
		} else {
			// Error during connect
			status = TargetStatus{Target: &t, Online: true, Since: time.Now()}
			conn.Close()
		}

		res <- status

		// waiting to ticker
		<-ticker
	}
	end <- 1
}
