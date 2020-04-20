package main

import (
 	"net"
	"fmt"
	"golang.org/x/net/icmp"
    	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)



func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func isIPv6(ip net.IP) bool {
	return len(ip) == net.IPv6len
}


func (p* Packet) makeECHOmessage() ([] byte , error) {
	body := &icmp.Echo{
		ID:   p.id,
		Seq:  p.sequence,
		Data: []byte(messageString),
	}
	var message icmp.Message
	if p.protocol=="ipv4"{
		message = icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Body: body,
	}
    	} else {
		message = icmp.Message{
		Type: ipv6.ICMPTypeEchoRequest,
		Code:0,
		Body: body,
		}
	     }
		
    b, err := message.Marshal(nil)
    if err != nil {
	fmt.Printf(" Error in Marshelling the Echo Packet : Error %s",err)
        return b, err
    }
    return b, err
}





func (p *Packet)resolveIPAddress(){
	switch p.protocol {
	  case "ipv4":
		p.IP, p.err = net.ResolveIPAddr("ip4", p.address)
	  case "ipv6":
		p.IP, p.err = net.ResolveIPAddr("ip6", p.address)
	  default:
		p.IP, p.err = net.ResolveIPAddr("ip", p.address)
		if (p.err!=nil){
                        p.IP.IP =nil
     		}
	}
	fmt.Printf("ICMP Echo Request to %s (%s) %d data Bytes\n",p.address,p.IP.IP,maxMessageLength)
}


func (p *Packet) returnParsedMessage(n int,reply []byte) (*icmp.Message , error){
       if p.protocol=="ipv4" {
	        rm, err := icmp.ParseMessage(ProtocolIvP4, reply[:n])
		return rm,err
	} else {	
		rm, err := icmp.ParseMessage(ProtocolIvP6, reply[:n])
		return rm,err
	}
}
