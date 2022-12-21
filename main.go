package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const help string = "Usage: tracert [-I] host"

func getStringIP(addr net.Addr) string {
	switch add := addr.(type) {
	case *net.UDPAddr:
		return add.IP.String()
	case *net.TCPAddr:
		return add.IP.String()
	case *net.IPAddr:
		return add.IP.String()
	}
	return ""
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var target = "ya.ru"
var maxHops = 64

func main() {
	icmp_style := false
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-I":
			icmp_style = true
			//make icmp type tracing
		case "--help", "-h":
			fmt.Println(help)
			os.Exit(0)
		default:
			target = os.Args[i]
		}
	}

	ip, err := net.ResolveIPAddr("ip4", target)

	// like "no such host", or something
	// Include check for invalid target
	catch(err)

	fmt.Printf("tracing route to %s (%s), %d hops max\n", target, ip.String(), maxHops)

	var network string
	var IP net.Addr

	if icmp_style {
		// icmp style
		network = "ip4:icmp"
		IP = ip
	} else {
		// udp style
		// You can think, what it sends udp packet, but actually it sends real ICMP packet
		// Without root!!!
		// (Really, there is no difference between icmp-style)
		network = "udp4"
		IP = &net.UDPAddr{IP: ip.IP}
	}
	conn, err := icmp.ListenPacket(network, "0.0.0.0")
	catch(err)
	defer conn.Close()

	var reply icmp.Type = ipv4.ICMPTypeTimeExceeded
	for hop := 1; reply != ipv4.ICMPTypeEchoReply && hop < maxHops; hop++ {
		reply = Ping(IP, conn, hop)
	}
}

func Ping(ip net.Addr, conn *icmp.PacketConn, ttl int) icmp.Type {
	conn.IPv4PacketConn().SetTTL(ttl)

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff,
			//ttl better than 1
			Seq:  ttl,
			Data: []byte(""),
		},
	}
	msg_bytes, _ := msg.Marshal(nil)

	// Write the message to the listening connection
	reply := make([]byte, 1500)

	conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	_, err := conn.WriteTo(msg_bytes, ip)
	// we can catch "no route to host" or something
	catch(err)
	// time for round-trip
	t1 := time.Now()

	n, add, err := conn.ReadFrom(reply)
	t2 := time.Now()

	if err != nil {
		fmt.Printf("%d * * *\n", ttl)
		return ipv4.ICMPTypeTimeExceeded
	} else {
		// 1 = iana.ProtocolICMP
		parsed_reply, err := icmp.ParseMessage(1, reply[:n])

		if err != nil {
			fmt.Printf("Error on ParseMessage %v\n", err)
			return ipv4.ICMPTypeTimeExceeded
		}
		// Resolution domain names
		// list of names for this address
		nameList, _ := net.LookupAddr(getStringIP(add))
		//round-trip time
		dt := t2.Sub(t1)
		fmt.Println(ttl, getStringIP(add), nameList, dt.String())

		switch parsed_reply.Type {
		case ipv4.ICMPTypeEchoReply:
			// Good
		case ipv4.ICMPTypeTimeExceeded:
			// expected
		case ipv4.ICMPTypeDestinationUnreachable:
			// ?????
			fmt.Printf("Host %s is unreachable\n", target)
		default:
			// We don't know what this is, so we can assume it's unreachable
			fmt.Printf("Host %s is unreachable\n", target)
		}
		return parsed_reply.Type
	}
}
