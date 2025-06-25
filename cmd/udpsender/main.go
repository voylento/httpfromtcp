package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const address = "localhost:42069"

func main() {

	udpaddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("udpaddr, IP=%v, Port=%d, Zone=%s\n", udpaddr.IP, udpaddr.Port, udpaddr.Zone)

	udpConn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading string from stdin: %v\n", err)
		}

		_, err = udpConn.Write([]byte(str))
		if err != nil {
			log.Print(err)
		}
	}


}

