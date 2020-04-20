package main

import (
	"fmt"
	"time"
	"net"
    	"golang.org/x/net/ipv6"
)


func (p* Packet) IPv6() {
	//fmt.Println("message",message)
	//fmt.Println("Protocol",p.protocol)
	//fmt.Println("Count to send",p.count)
	//fmt.Println("ttl",p.ttl)
	//fmt.Println("address",p.address)
	//fmt.Println("delay",p.delay)

	conn, err := net.ListenPacket("ip6:58", ListenAddrIPv6)
	if err != nil {
		fmt.Printf("Unable to create a connection object : Error %s",err)
		return 
	}
	//fmt.Println("Connection obj,protocol,err",conn,p.protocol,p.err)

	p.resolveIPAddress()
        //fmt.Println("Resolved address",p.IP)
	if p.IP.IP == nil {
		//fmt.Println("Unable to resolve IP address")
		return 
	}

	defer conn.Close()  // Close the connection when everything is done

	packetipv6 := ipv6.NewPacketConn(conn)
	
	if err := packetipv6.SetControlMessage(ipv6.FlagHopLimit|ipv6.FlagDst|ipv6.FlagInterface, true); err != nil {
		fmt.Printf("Unable to set Control message for IPv6 protocol : Error %s",err)
	}


	var f ipv6.ICMPFilter
	f.SetAll(true)
	f.Accept(ipv6.ICMPTypeTimeExceeded)
	f.Accept(ipv6.ICMPTypeEchoReply)

	if err := packetipv6.SetICMPFilter(&f); err != nil {
		fmt.Printf("Unable to set filter for IPv6 protocol : Error %s",err)
	}
	var wcm ipv6.ControlMessage
	reply := make([]byte, ReadDataSize)
	p.sequence += 1
    	for { 	// Run Infinite Loop for the code
		
		if p.count!=0 && p.sequence==p.count{
			break
		}
		b,err := p.makeECHOmessage()
		if err != nil{
			return
		}
		p.sequence += 1
		wcm.HopLimit = p.ttl
	

		//fmt.Println("Time in second before sending ",time.Now())
		time.Sleep(time.Duration(p.delay)*time.Second)
		start := time.Now()
		n, err := packetipv6.WriteTo(b,&wcm,p.IP)
		if (err!=nil){
			fmt.Printf("error to write data : Error %s\n",err)
			return 
		}

		
		//fmt.Println("Time in second after adding ",time.Now().Add(time.Duration(p.delay)*time.Second))

		err = conn.SetReadDeadline(time.Now().Add(time.Duration(ReadDeadline)*time.Second)) // Read deadline is 5 seconds
		if err != nil {
			fmt.Printf("Unable to set Read Deadline : Error %s\n",err)
			return 
		}
		n,cm, _ , err := packetipv6.ReadFrom(reply)
		if err != nil {
			fmt.Printf("%d bytes received from %s target (%s): Loss\n",n,p.address,cm)
			continue
		}

		duration := time.Since(start)
		//fmt.Println("Error due to read failure",err)

		rm, err := p.returnParsedMessage(n,reply)
		Seq:= rm.Body.(*icmp.Echo).Seq 
		
		if p.ttl<cm.HopLimit {
				fmt.Printf("From %s (%s) icmp_seq=%d Time exceeded: Hop limit \n",cm.Dst,cm.Dst,p.sequence)
				continue
			}
				
		switch rm.Type {
		   case ipv6.ICMPTypeTimeExceeded:
			fmt.Printf("%d bytes recieved from %s target (%s): icmp_seq=%d time=%s : Loss\n",n,p.address,p.IP,p.sequence,duration)
		   case ipv6.ICMPTypeEchoReply:
			fmt.Printf("%d bytes from %s (%s): icmp_seq=%d ttl=%d time=%.4s ms \n",n,p.address,p.IP.IP,Seq,cm.HopLimit,duration)
			
		}
		
	    }
}



