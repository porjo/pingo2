package main

import (
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
)

const ICMPReadTimeout = 2
const ICMPWriteTimeout = 2

// non-privileged ping on Linux requires special sysctl setting:
//     sysctl -w net.ipv4.ping_group_range="0 0"
//
// Where group matches running process
// See: http://stackoverflow.com/questions/8290046/icmp-sockets-linux/20105379#20105379
func Ping(hostname string) (reply bool, err error) {
	ipAddr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		return false, err
	}

	readDeadline := time.Now().Add(time.Duration(time.Second * ICMPReadTimeout))
	writeDeadline := time.Now().Add(time.Duration(time.Second * ICMPWriteTimeout))

	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return false, err
	}
	defer c.Close()

	if err = c.SetReadDeadline(readDeadline); err != nil {
		return false, err
	}
	if err = c.SetWriteDeadline(writeDeadline); err != nil {
		return false, err
	}

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return false, err
	}
	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: ipAddr.IP}); err != nil {
		log.Fatalf("WriteTo err, %s", err)
	}

	rb := make([]byte, 1500)
	n, _, err := c.ReadFrom(rb)
	if err != nil {
		return false, err
	}
	rm, err := icmp.ParseMessage(iana.ProtocolICMP, rb[:n])
	if err != nil {
		return false, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return true, nil
	default:
		return false, nil
	}

	return false, nil
}
