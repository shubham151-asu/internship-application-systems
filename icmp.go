package main

import (
	"fmt"
	"time"
    	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)


func (p* Packet) ICMP() {
	//fmt.Println("Protocol",p.protocol)
	//fmt.Println("Count to send",p.count)
	//fmt.Println("ttl",p.ttl)
	//fmt.Println("address",p.address)
	//fmt.Println("delay",p.delay)

	conn, err := icmp.ListenPacket("ip4:icmp", ListenAddrICMP)
	if err != nil {
		fmt.Println("Unable to create a connection object")
		return
	}
	//fmt.Println("Connection obj,protocol,err",conn,p.protocol,p.err)

	//p.resolveIPAddress()
        //fmt.Println("Resolved address",p.IP.IP)
	if p.IP == nil {
		fmt.Printf("Unable to resolve IP address (%s) \n",p.address)
		return 
	}

	defer conn.Close()  // Close the connection when everything is done

	reply := make([]byte, ReadDataSize)
	p.sequence += 1
    	for { 	// Run Infinite Loop for the code

		if p.count!=0 && p.sequence==p.count{
			break
		}
		b,err := p.makeECHOmessage()
		if err != nil{
			fmt.Printf("Unable to make Echo Message : Error %s\n",err)
			return 
		}
		p.sequence += 1

		time.Sleep(time.Duration(p.delay)*time.Second)
		start := time.Now()
		n, err := conn.WriteTo(b, p.IP)

		
		//fmt.Println("Time in second after adding ",time.Now().Add(cmdArg.delay* time.Second))

		err = conn.SetReadDeadline(time.Now().Add(time.Duration(ReadDeadline)*time.Second)) // Read deadline is 5 seconds
		if err != nil {
			fmt.Println("Unable to set Read Deadline")
			return
		}
		n, _ , err = conn.ReadFrom(reply)

		duration := time.Since(start)
		//fmt.Println("Error due to read failure",err)
		if err != nil {
			fmt.Printf("%d bytes received from %s target (%s): Loss\n",n,p.address,p.IP.IP)
			continue
		}

		rm, err := p.returnParsedMessage(n,reply)
	
		if (p.protocol=="ipv4" && rm.Type!= ipv4.ICMPTypeEchoReply){
			fmt.Printf("%d bytes received from %s target (%s): Loss\n",n,p.address,p.IP.IP)
                	continue
		}

		var Seq int	
		if p.protocol=="ipv4" {
			Seq = rm.Body.(*icmp.Echo).Seq
		}

		fmt.Printf("%d bytes from %s (%s): icmp_seq=%d time=%.4s ms \n",n,p.address,p.IP.IP,Seq,duration)
		if err != nil {
			return 
		}
		//fmt.Println("Duration and peer:",peer,duration)
		
		if p.count!=0 && Seq==p.count{
			break
		}
	    }
   
}



