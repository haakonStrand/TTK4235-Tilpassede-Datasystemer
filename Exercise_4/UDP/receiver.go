package main

import (
	"fmt"
	"net"
)

func main() {
	// Resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", "localhost:8888")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// Listen for UDP packets
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	// Buffer to hold received data
	buffer := make([]byte, 1024)

	// Receive data
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	fmt.Println("Received message:", string(buffer[:n]))
}
