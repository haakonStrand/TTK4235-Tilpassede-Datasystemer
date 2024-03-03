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

	// Establish UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Message to send
	message := []byte("Hello from sender!")

	// Send message
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Println("Message sent successfully")
}
