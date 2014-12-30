package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const CheckInterval = 30

type Target struct {
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
	log.Println("starting runtarget on ", t.Name)
	if t.Interval < CheckInterval {
		t.Interval = CheckInterval
	}
	ticker := time.Tick(time.Duration(t.Interval) * time.Second)
	status := TargetStatus{Target: &t, Online: true, Since: time.Now()}
	for {
		var err error
		var conn net.Conn
		var failed bool
		var addrURL *url.URL

		status.ErrorMsg = ""

		addrURL, err = url.Parse(t.Addr)
		if err != nil {
			log.Printf("target address %s could not be read, %s", t.Addr, err)
			break
		}

		// Polling
		if addrURL.Scheme == "http" || addrURL.Scheme == "https" {
			var resp *http.Response
			var client *http.Client

			req, _ := http.NewRequest("GET", addrURL.String(), nil)
			if t.Host != "" {
				// Set hostname for TLS connection. This allows us to connect to use
				// another hostname or IP for the actual TCP connection. Handy for GeoDNS scenarios.
				transport := http.Transport{
					TLSClientConfig: &tls.Config{
						ServerName: t.Host,
					},
				}
				req.Host = t.Host
				client = &http.Client{Transport: &transport}
			} else {
				client = &http.Client{}
			}
			resp, err = client.Do(req)
			if err != nil {
				log.Printf("Error %s\n", err)
				status.ErrorMsg = fmt.Sprintf("%s", err)
				failed = true
			} else {
				var body []byte
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error %s\n", err)
					status.ErrorMsg = fmt.Sprintf("%s", err)
					failed = true
				} else {
					resp.Body.Close()

					if t.Keyword != "" {
						if strings.Index(string(body), t.Keyword) == -1 {
							status.ErrorMsg = fmt.Sprintf("keyword '%s' not found", t.Keyword)
							log.Printf("%s, %s\n", t.Name, status.ErrorMsg)
							failed = true
						}
					}
				}
			}
		} else {
			conn, err = net.DialTimeout("tcp", addrURL.Host, time.Duration(config.Timeout)*time.Second)
			if err != nil {
				log.Printf("Error %s\n", err)
				status.ErrorMsg = fmt.Sprintf("%s", err)
				failed = true
			} else {
				conn.Close()
			}
		}

		if failed {
			// Error during connect
			status.Online = false
			if time.Since(status.LastAlert) > time.Second*time.Duration(config.Alert.Interval) {
				alert(&status, config)
			}
			if status.Online {
				status.Since = time.Now()
			}
		} else {
			// Connect ok
			if !status.Online {
				status.Since = time.Now()
				status.Online = true
				alert(&status, config)
			} else {
				status.Online = true
			}
		}
		status.LastCheck = time.Now()

		res <- status

		// waiting for ticker
		<-ticker
	}
}

func alert(status *TargetStatus, config Config) {
	if config.Alert.ToEmail != "" {
		err := EmailAlert(*status, config)
		if err != nil {
			log.Printf("%s\n", err)
		}
		status.LastAlert = time.Now()
	}
}
