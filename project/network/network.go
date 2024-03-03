package network

import (
	"Project/network/network_driver/bcast"
	"Project/network/network_driver/peers"
	"fmt"
	"time"
)

func InitNetwork(chanReceive chan []byte, chanMsgToBeSent chan []byte, id int) {
    
    peerUpdateCh := make(chan peers.PeerUpdate)
    peerTxEnable := make(chan bool)
    go peers.Transmitter(30003, fmt.Sprint(id), peerTxEnable)
    go peers.Receiver(30003, peerUpdateCh)

    networkTx := make(chan []byte)
    networkRx := make(chan []byte)
    go bcast.Transmitter(20003, networkTx)
    go bcast.Receiver(20003, networkRx)

    go func() {
        ticker := time.NewTicker(1 * time.Second)//10 * time.Millisecond
        defer ticker.Stop()
        lastMsg := []byte("")
        for {
            select {
            case <-ticker.C:
                networkTx <- lastMsg
            case sendMsg := <-chanMsgToBeSent:
                networkTx <- sendMsg
                lastMsg = sendMsg
            }
        }
    }()

    for {
        select {
        case p := <-peerUpdateCh:
            //ONLY FOR DEBUGGING
            fmt.Printf("Peer update:\n")
            fmt.Printf(" Peers: %q\n", p.Peers)
            fmt.Printf(" New: %q\n", p.New)
            fmt.Printf(" Lost: %q\n", p.Lost)

        case receiveMsg := <-networkRx:
            chanReceive <- receiveMsg
        }
    }
}