package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	ADD := "127.0.0.1:42069"
	// nc -u -l 42069
	// use this command to listen for udp in 42069
	udp, err := net.ResolveUDPAddr("udp", ADD)

	if err != nil {
		fmt.Printf("Failed to resolve UDP eddr with err %v", err)
		panic("Failed to resolve UDP Add")
	}

	// radd, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")

	cn, _ := net.DialUDP("udp", nil, udp)
	defer func() {
		fmt.Println(cn)
		err := cn.Close()
		fmt.Printf("Failed to close connections %v", err)
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("Error reading from the line %v ]n", err)
			continue
		}

		cn.Write(line)
	}

}
