package main

import (
	"fmt"
	"net"
	"time"
)

func listenForMessage() {
	address, err := net.ResolveUDPAddr("udp", ":20003")
	if err != nil {
		fmt.Println("There was an error resolving address")
	}
	ln, err := net.ListenUDP("udp", address)
	if err != nil {
		fmt.Println("There was an error")
	}

	for {
		buffer := make([]byte, 1024)
		n, addr, err := ln.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("There was an error reading")
		}
		fmt.Println("Received message from", addr, ": ", string(buffer[:n]))

	}

}
func sendMessage() {
	address, err := net.ResolveUDPAddr("udp", "10.100.23.255:20003")
	if err != nil {
		fmt.Println("There was an error resolving address")
	}
	ln, err := net.DialUDP("udp", nil, address)
	if err != nil {
		fmt.Println("There was an error")
	}

	for {
		time.Sleep(500 * time.Millisecond)
		_, err := ln.Write([]byte("Hello, this is group 32 broadcasting"))
		if err != nil {
			fmt.Println("There was another different error in the sendingMessage rutine")
		}
	}

}

// func sendMessage() {

// 	ln, err := net.ListenPacket("udp4",":20003")
// 	if err != nil{
// 		fmt.Println("There was an error in the sendMessage rutine")

// 	}
// 	defer ln.Close()

// 	addr, err := net.ResolveUDPAddr("udp4", "10.100.23.255:20003")
// 	if err != nil{
// 		fmt.Println("There was a different error in the sendingMessage rutine")
// 	}

// 	for {
// 	time.Sleep(500*time.Millisecond)
// 	_ , err := ln.WriteTo([]byte("Hello, this is group 32 broadcasting"), addr)
// 	if err != nil{
// 		fmt.Println("There was another different error in the sendingMessage rutine")
// 	}
// 	}
// }

func main() {

	go listenForMessage()
	go sendMessage()
	select {}
}

//IP-address: 10.100.23.129
