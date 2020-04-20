package main

import (
	"fmt"
	"time"
	"net"
    	"golang.org/x/net/icmp"
    	"golang.org/x/net/ipv4"
)


func (p* Packet) IPv4() {
	//fmt.Println("Protocol",p.protocol)
	//fmt.Println("Count to send",p.count)
	//fmt.Println("ttl",p.ttl)
	//fmt.Println("address",p.address)
	//fmt.Println("delay",p.delay)

	conn, err := net.ListenPacket("ip4:1", ListenAddrIPv4)
	if err != nil {
		fmt.Printf("Unable to create a connection object : Error %s\n",err)
		return 
	}
	//fmt.Println("Connection obj,protocol,err",conn,p.protocol,p.err)

	p.resolveIPAddress()
        //fmt.Println("Resolved address",p.IP)
	if p.IP.IP == nil  {
		fmt.Printf("Unable to resolve IP address %s\n",p.address)
		return 
	}

	defer conn.Close()  // Close the connection when everything is done

	packetipv4 := ipv4.NewPacketConn(conn)
	
	if err := packetipv4.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		fmt.Println("Unable to set Control message in IPv4")
		return
	}

	//packetipv4.SetTOS(0x0)
	packetipv4.SetTTL(p.ttl)
	reply := make([]byte, ReadDataSize)
	//count := 1
	p.sequence += 1
    	for { 	// Run Infinite Loop for the code
		
		if p.count!=0 && p.sequence==p.count{
			break
		}

		b,err := p.makeECHOmessage()
		if err != nil{
			fmt.Println("Unable to make ICMP Echo message : Error",err)
			return 
		}
		p.sequence += 1
		

		//fmt.Println("Time in second before sending ",time.Now())
		time.Sleep(time.Duration(p.delay)*time.Second)

		start := time.Now()
		
		n, err := packetipv4.WriteTo(b,nil,p.IP)
		if (err!=nil){
			fmt.Println("error to write data using IPv4 protocol : Error",err)
			return 
		}

		err = conn.SetReadDeadline(time.Now().Add(time.Duration(ReadDeadline)*time.Second)) // Read deadline is 5 seconds
		if err != nil {
			fmt.Printf("Unable to set Read Deadline : Error %s",err)
			return
		}
		n,cm, _ , err := packetipv4.ReadFrom(reply)

		duration := time.Since(start)
		//fmt.Println("details from cm and peer",cm.TTL,cm.,peer)
		rm, err := p.returnParsedMessage(n,reply)

		if (rm.Type!= ipv4.ICMPTypeEchoReply){
			fmt.Printf("%d bytes received from %s target (%s): Loss\n",n,p.address,p.IP.IP)
                	continue
		}
		
		Seq:= rm.Body.(*icmp.Echo).Seq  // Check the Sequence of the message
		
		if p.ttl<cm.TTL {
				fmt.Printf("From %s (%s) icmp_seq=%d Time to live exceeded \n",cm.Dst,cm.Dst,p.sequence)
				continue
			}
		
		if err != nil {
			fmt.Printf("%d bytes received from %s target (%s) : Loss : Error %s\n",n,p.address,p.IP.IP,err)
			continue
		}

		

		fmt.Printf("%d bytes from %s (%s): icmp_seq=%d ttl=%d time=%.4s ms \n",n,p.address,p.IP.IP,Seq,cm.TTL,duration)
		if err != nil {
			return
		}
		//fmt.Println("Duration and peer:",peer,duration)
		
	    }
}



