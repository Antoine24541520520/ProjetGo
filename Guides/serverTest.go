// go build serverTest.go && ./serverTest
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

var (
	clientCount int
	clients     = make(map[net.Conn]struct{})
	mu          sync.Mutex
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		mu.Lock()
		if clientCount >= maxClients {
			mu.Unlock()
			continue
		}

		clientCount++
		clients[conn] = struct{}{}
		if clientCount == maxClients {
			for clientConn := range clients {
				fmt.Fprintln(clientConn, "Ready")
			}
		}
		mu.Unlock()

		wg.Add(1)
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

		if msg == "space\n" {
			fmt.Printf("%s a appuyé sur Espace \n", conn.RemoteAddr().String())
			mu.Lock()
			for clientConn := range clients {
				if clientConn != conn {
					fmt.Fprintln(clientConn, conn.RemoteAddr().String()+" a appuyé sur Espace")
				}
			}
			mu.Unlock()
		}

		if _, err := writer.WriteString("Received: " + msg); err != nil {
			fmt.Printf("Error writing message to %s: %v\n", conn.RemoteAddr(), err)
			break
		}

		if err := writer.Flush(); err != nil {
			fmt.Printf("Error flushing message to %s: %v\n", conn.RemoteAddr(), err)
			break
		}
	}

	mu.Lock()
	clientCount--
	delete(clients, conn)
	mu.Unlock()

	fmt.Printf("Client %s disconnected\n", conn.RemoteAddr())
}
