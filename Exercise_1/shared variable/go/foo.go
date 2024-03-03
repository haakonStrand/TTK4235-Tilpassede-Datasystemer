// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
)

var i = 0

func incrementing(chan1 chan int, chan3 chan int) {
	//TODO: increment i 1000000 times
	for j := 0; j < 999998; j++ {
		chan1 <- 1
	}
	chan3 <- 1
}

func decrementing(chan2 chan int, chan3 chan int) {
	//TODO: decrement i 1000000 times
	for j := 0; j < 1000000; j++ {
		chan2 <- 1
	}
	chan3 <- 1
}

func server(chan1 chan int, chan2 chan int, chan3 chan int, read chan int) {
	for {
		select {
		case <-chan1:
			i++
		case <-chan2:
			i--
		case read <- i:
		}
	}
}

func main() {
	//GOMAXPROCS limits the number of operating system threads that can execute Go code simultaneously
	//Set it to 1 - only one thread can execute, so the number will either increment og decrement continously.
	runtime.GOMAXPROCS(2)

	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)
	read := make(chan int)

	// TODO: Spawn both functions as goroutines
	go incrementing(chan1, chan3)
	go decrementing(chan2, chan3)
	go server(chan1, chan2, chan3, read)

	<-chan3
	<-chan3

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	Println("The magic number is:", <-read)

}
