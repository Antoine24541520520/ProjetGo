package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	maxClients = 4
)

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Printf("Error listening on port 1234: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on localhost:1234...")

	var wg sync.WaitGroup
	clientCount := 0
	for {
		if clientCount >= maxClients {
			break
		}

		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		wg.Add(1)
		clientCount++
		go handleConnection(conn, &wg)
	}

	wg.Wait()
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	fmt.Printf("Client %s connected\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading message from %s: %v\n", conn.RemoteAddr(), err)
			}
			break
		}

		fmt.Printf("Message from %s: %s", conn.RemoteAddr(), msg)

		if _, err := writer.WriteString("Received: " + msg); err != nil {
			fmt.Printf("Error writing message to %s: %v\n", conn.RemoteAddr(), err)
			break
		}

		if err := writer.Flush(); err != nil {
			fmt.Printf("Error flushing message to %s: %v\n", conn.RemoteAddr(), err)
			break
		}
	}

	fmt.Printf("Client %s disconnected\n", conn.RemoteAddr())
}
