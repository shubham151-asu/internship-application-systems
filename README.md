Description
---------------
Code to build a small CLI application using golang to send icmp "echo requests" to hosts and receive "echo reply" from the host. Code takes hostname or IP from CLI and sends ICMP request to the host in an infinite loop. Code also allows to set TTL limits, packets count,interval params from the CLI similar to ping.


Current Implementation and TODO items
---------------
	1. Able to send ICMP Echo request and recieve Echo Reply using icmp protocol
	2. Able to send ICMP Echo request and recieve Echo Reply using ipv4 protocol
	3. Able to send ICMP Echo request and recieve Echo Reply using ipv6 protocol
	4. Able to set TTL limits for IPv4 and HopLimits in IPv6 
	5. Set logs in the code : TODO
	6. Ping statistics : TODO
	
	
Dependencies
---------------
	1.Go compiler
	2.Packages :  net
	             "golang.org/x/net/icmp"
		     "golang.org/x/net/ipv4"
		     "golang.org/x/net/ipv6" 


Build
----------------

go build main.go icmp.go ipv4.go ipv6.go constants.go util.go

Run
----------------
./main params     ( if unable to create connection object : use sudo ./main params )

	params
	----------
	Mandatory params :
		-h : to provide hostname or IP 
	Optional params :
		-t : to set TTL limit (an integer)
		-c : to set number of packets to be sent (an integer)
		-i : Interval between each package sent (an integer)
		-p : protocol : can be icmp,ipv4,ipv6 (Default protocol is IPv4 )
				with icmp ttl/hopLimits will not be set

	sample run :
		./main -h=google.com -p=ipv4 -t=45 -c=65 -i=2 
		This will send ICMP Echo message to google IPv4 address after resolution with TTL limit 45,
		total number of packets to be sent is 65, interval or delay between each message is 2 seconds










