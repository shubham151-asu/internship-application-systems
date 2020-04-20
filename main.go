package main

import (
	"fmt"
 	"net"
	"flag"
)


type Packet struct {
	protocol string
	count int
	ttl int
	sequence int
	id int
	IP *net.IPAddr
	address string
	delay int
	tos int
	err error
}

type ICMPStatistic struct {
	duration float64
	transmitted int
	received int
	Loss int
}

func main() {

	protocol := flag.String("p","ipv6","Default protocol is IPv4")
	count := flag.Int("c",DefaultCount,"Default number of packets to be sent is 64")
	ttl := flag.Int("t",DefaultTTL,"Default TTL is 59")
	host := flag.String("h","","Default host is not set")
	delay := flag.Int("i",DefaultDelay,"Default delay in Packet sent is 1s")
	flag.Parse()
	if *host==""{
		fmt.Println("Host not assigned : Please specify host")
		return
	}

	p:= Packet{
		protocol: *protocol,
		count:*count,
		ttl:*ttl,
		address:*host,
		delay:*delay, 
	}
	

	switch *protocol {
		case "icmp":
			p.protocol = "ipv4"
			p.ICMP()
		case "ipv4":
			p.protocol = "ipv4"
			p.IPv4()
		case "ipv6":
			p.protocol = "ipv6"
			p.IPv6()
	}
	
}
