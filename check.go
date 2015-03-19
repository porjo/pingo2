package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// minimum interval between checks. Used as default value when none set by user.
const CheckInterval = 30

type Target struct {
	// target id
	Id int
	// Name of the Target
	Name string
	// Address of the target e.g. "http://localhost"
	Addr string
	// HTTP 'Host:' header (if different from Addr)
	Host string
	// Polling interval, in seconds
	Interval int
	// Look for this string in the response body
	Keyword string
}

type TargetStatus struct {
	Target    *Target
	Online    bool
	ErrorMsg  string
	Since     time.Time
	LastCheck time.Time
	LastAlert time.Time
}

func startTarget(t Target, res chan TargetStatus, config Config) {
	go runTarget(t, res, config)
}

func runTarget(t Target, res chan TargetStatus, config Config) {
	var err error
	var failed bool
	var addrURL *url.URL
	log.Printf("starting runtarget on %s", t.Name)
	if t.Interval < CheckInterval {
		t.Interval = CheckInterval
	}

	// wait a bit, to randomize check offset
	time.Sleep(time.Duration(rand.Intn(t.Interval)) * time.Second)

	ticker := time.Tick(time.Duration(t.Interval) * time.Second)
	alertRequest := make(chan *TargetStatus, 1)
	// spawn routine to handle alert requests
	go alertRoutine(alertRequest, config)
	status := TargetStatus{Target: &t, Online: true, Since: time.Now()}

	for {
		failed = false
		status.ErrorMsg = ""

		addrURL, err = url.Parse(t.Addr)
		if err != nil {
			log.Printf("target address %s could not be read, %s", t.Addr, err)
			break
		}

		// Polling
		switch addrURL.Scheme {
		case "http", "https":
			var resp *http.Response
			var client *http.Client

			req, _ := http.NewRequest("GET", addrURL.String(), nil)
			transport := &http.Transport{
				DisableKeepAlives:  true,
				DisableCompression: true,
			}
			if t.Host != "" {
				// Set hostname for TLS connection. This allows us to connect using
				// another hostname or IP for the actual TCP connection. Handy for GeoDNS scenarios.
				transport.TLSClientConfig = &tls.Config{
					ServerName: t.Host,
				}
				req.Host = t.Host
			}
			client = &http.Client{
				Timeout:   time.Duration(config.Timeout) * time.Second,
				Transport: transport,
			}
			resp, err = client.Do(req)
			if err != nil {
				log.Printf("[%d:%s] http(s) error, %s", t.Id, addrURL, err)
				status.ErrorMsg = fmt.Sprintf("%s", err)
				failed = true
			} else {
				var body []byte
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("[%d:%s] http(s) error, %s", t.Id, addrURL, err)
					status.ErrorMsg = fmt.Sprintf("%s", err)
					failed = true
				} else {
					if t.Keyword != "" {
						if strings.Index(string(body), t.Keyword) == -1 {
							status.ErrorMsg = fmt.Sprintf("keyword '%s' not found", t.Keyword)
							log.Printf("[%d:%s] http(s) error, %s", t.Id, addrURL, status.ErrorMsg)
							failed = true
						}
					}
				}
				resp.Body.Close()
			}
		case "ping":
			var success bool
			success, err = Ping(addrURL.Host)
			if err != nil {
				log.Printf("[%d:%s] ping error, %s", t.Id, addrURL, err)
				status.ErrorMsg = fmt.Sprintf("%s", err)
			}
			failed = !success
		default:
			var conn net.Conn
			conn, err = net.DialTimeout("tcp", addrURL.Host, time.Duration(config.Timeout)*time.Second)
			if err != nil {
				log.Printf("[%d:%s] tcp conn error, %s", t.Id, addrURL, err)
				status.ErrorMsg = fmt.Sprintf("%s", err)
				failed = true
			} else {
				conn.Close()
			}
		}

		status.LastCheck = time.Now()

		if debug {
			log.Printf("[%d:%s] failed=%v, online=%v, since=%s, last_alert=%s, last_check=%s", t.Id, addrURL, failed, status.Online, status.Since, status.LastAlert, status.LastCheck)
		}

		if failed {
			// Error during connect
			if status.Online {
				// was online, now offline
				status.Online = false
				status.Since = time.Now()
				alertRequest <- &status

			} else {
				// was offline, still offline
				if time.Since(status.LastAlert) > time.Second*time.Duration(config.Alert.Interval) {
					alertRequest <- &status
				}
			}
		} else {
			// Connect ok
			if !status.Online {
				// was offline, now online
				status.Online = true
				lastSince := status.Since
				status.Since = time.Now()
				if debug {
					log.Printf("[%d:%s] was offline, now online - time since=%s", t.Id, addrURL, time.Since(lastSince))
				}
				// Don't bother with 'up' alert if the host was down less than a minute
				if time.Since(lastSince) > time.Duration(time.Minute) {
					alertRequest <- &status
				}
			}
		}

		res <- status

		// waiting for ticker
		<-ticker
	}
}

func alert(status *TargetStatus, config Config) {
	if config.Alert.ToEmail != "" {
		err := EmailAlert(*status, config)
		if err != nil {
			log.Printf("%s", err)
		}
		if debug {
			log.Printf("[%d:%s] alert sent to %s", status.Target.Id, status.Target.Addr, config.Alert.ToEmail)
		}
	} else {
		if debug {
			log.Printf("[%d:%s] alert NOT sent as no 'To:' email specified", status.Target.Id, status.Target.Addr)
		}
	}
	status.LastAlert = time.Now()
}

func alertRoutine(alertRequest <-chan *TargetStatus, config Config) {
	for {
		select {
		case req := <-alertRequest:
			// Host is online, or has been offline for greater than a minute
			if req.Online || time.Since(req.Since) > time.Duration(time.Minute) {
				alert(req, config)
			} else {
				// Don't bother with 'down' alert
				// if host comes back within a minute
				timer1 := time.NewTimer(time.Minute)
				for {
					select {
					case req2 := <-alertRequest:
						if req2.Online {
							alert(req2, config)
							goto done
						} else {
							// if another 'offline' requests comes in the meantime
							// ignore and continue waiting for timer or 'online'
							continue
						}
					case <-timer1.C:
						alert(req, config)
					}

				done:
					break
				}
			}
		}
	}
}
