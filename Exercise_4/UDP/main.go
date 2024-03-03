package main

import (
	"fmt"
	"net"
	"time"
	"os/exec"
	"strconv"
)

func server() int{
	
		// Resolve UDP address
		addr, err := net.ResolveUDPAddr("udp", "localhost:8888")
		if err != nil {
			fmt.Println("Error resolving address:", err)
			return 0
		}
	
		// Listen for UDP packets
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("Error listening:", err)
			return 0
		}
		defer conn.Close()

		lastValue := 0
	
		for {
		
		conn.SetDeadline(time.Now().Add(3 * time.Second))
	
		// Buffer to hold received data
		buffer := make([]byte, 1024)
	
		// Receive data
		n, _, err := conn.ReadFromUDP(buffer)
		
		
		
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Timeout error:", err)
				//close socket
				conn.Close()
				
				return lastValue
				} else if err != nil {
					fmt.Println("Error: ", err)
					return 0
				}
		}
		
		fmt.Println("Received message:", string(buffer[:n]))
		lastValue, _ = strconv.Atoi(string(buffer[:n]))
	}
	
}

func client(i int) {
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
	for {
		message := []byte(fmt.Sprint(i))

		// Send message
		_, err = conn.Write(message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
		i=i+1
		time.Sleep(1 * time.Second)
	}
}
	

func process() {
	lastValue := server()
	fmt.Println("Last value: ", lastValue)
	cmd := exec.Command("gnome-terminal", "--", "go", "run", "main.go")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	client(lastValue)
}

func main() {
	go process()
	select {}
}
