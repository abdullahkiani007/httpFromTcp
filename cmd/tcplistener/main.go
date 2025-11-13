package main

import (
	"fmt"
	"net"

	"github.com/abdullahkiani007/httpfromtcp/internal/request"
)

func main() {

	tcp, err := net.Listen("tcp", "127.0.0.1:42069")
	fmt.Printf("Listening to TCP\n")
	if err != nil {
		fmt.Printf("Error creating tcp listenere")
		panic("Error creating listerner\n")
	}

	for {
		fmt.Printf("\nListening for connections\n")
		cn, err := tcp.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection%s\n", err)
		}

		req, err := request.RequestFromReader(cn)
		rl := req.RequestLine
		if err != nil {
			fmt.Printf("Failed to parse error %e", err)
		}
		fmt.Println(req.Headers.Header)
		fmt.Printf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n", rl.Method, rl.RequestTarget, rl.HttpVersion)
		fmt.Printf("Headers:\n")
		for key, v := range req.Headers.Header {
			fmt.Printf("- %v: %v\n", key, v)
		}
	}

}
