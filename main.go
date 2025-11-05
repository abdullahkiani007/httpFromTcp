package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const inputFilePath = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {
	// var str string
	// var line string
	lb := []byte{}
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		for {
			b := make([]byte, 8)
			_, err := f.Read(b)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				break
			}
			// str = string(b)
			for _, v := range b {
				if v != 10 {
					lb = append(lb, v)
					// fmt.Printf("line _> %s\n", line)
				} else {

					// fmt.Printf("read: %s\n", string(lb))
					ch <- string(lb)
					lb = []byte{}
				}
			}

		}
		// fmt.Printf("read: %s\n", string(lb))
		ch <- string(lb)
	}()
	return ch
}

func main() {

	tcp, err := net.Listen("tcp", "127.0.0.1:42069")
	fmt.Printf("Listening to TCP\n")
	if err != nil {
		fmt.Printf("Error creating tcp listenere")
		panic("Error creating listerner\n")
	}

	for {
		fmt.Printf("Listening for connections\n")
		cn, err := tcp.Accept()
		if err != nil {
			fmt.Printf("could not open %s: %s\n", inputFilePath, err)
		}

		ch := getLinesChannel(cn)
		cn.Write([]byte("heyy"))
		for str := range ch {
			fmt.Printf("read: %v\n", str)
		}
	}

}
