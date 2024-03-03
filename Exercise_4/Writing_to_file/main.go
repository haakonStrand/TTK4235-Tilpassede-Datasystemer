package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func printAlive() {
	i, _ := os.ReadFile("backup.txt")
	if i == nil{
		i = make([]byte, 1)
	}
	i_int, _ := strconv.Atoi(string(i))
	for {
		i_int = i_int + 1
		i_byte := []byte(strconv.Itoa(i_int))
		os.WriteFile("backup.txt", i_byte, 0644)
		fmt.Println(i_int)	
		time.Sleep(1 * time.Second)	
	}
}

func primary() {
	//Run backupile, err := os.Open("backup.txt")
	time.Sleep(1 * time.Second)
	cmd := exec.Command("gnome-terminal", "--", "go", "run", "main.go")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	printAlive()
}

func backup() {
	//Check if primary is alive and become primary if not
	for {
		file, err := os.Open("backup.txt")
		if err != nil {
			panic(err)
		}
		fi, err := file.Stat()
		if err != nil {
			panic(err)
		}
		if fi.ModTime().Before(time.Now().Add(-3 * time.Second)) {
			go primary()
			return
		}
		time.Sleep(3 * time.Second)
	}

}

func main() {
	
	go backup()
	select {}
}
