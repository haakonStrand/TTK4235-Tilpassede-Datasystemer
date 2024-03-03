package main

import (
	"fmt"
	"net"
)


func listenForMessage(conn net.Conn){
	buffer := make([]byte, 1024)
	for{
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error 2")
		}
	
		fmt.Println("Message recieved: ",string( buffer[:n]))
	}

}


func main() {


	address, err := net.ResolveTCPAddr("tcp", "10.100.23.129:34933")
	conn, err := net.DialTCP("tcp",nil, address)
	if err != nil{
		fmt.Println("Error 1")
	}
	defer conn.Close()

	message := append([]byte("This is group 32 sending a message to the server "), 0)
	_ , err = conn.Write(message)
	if err != nil{
		fmt.Println("Error 2")
	}
	
	go listenForMessage(conn)



	select {}
}

//IP-address: 10.100.23.129
