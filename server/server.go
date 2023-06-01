package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	maxClients = 2
)

var (
	clientCount    int
	clients        = make(map[net.Conn]struct{})
	mu             sync.Mutex
	readyClients   int
	lockReady      bool
	lockReadyMutex sync.Mutex
	clientLocks    = make(map[net.Conn]bool)
	clientLocksMu  sync.Mutex
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
		clientLocksMu.Lock()
		clientLocks[conn] = false
		clientLocksMu.Unlock()
		fmt.Printf("num_client#%v", clientCount)
		for clientConn := range clients {

			fmt.Fprintf(clientConn, "num_client#%v\n", clientCount)
		}
		if clientCount == maxClients {
			lockReadyMutex.Lock()
			lockReady = true
			lockReadyMutex.Unlock()
			for clientConn := range clients {
				fmt.Fprintln(clientConn, "ready")
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
		if msg == "locked\n" {
			clientLocksMu.Lock()
			if !clientLocks[conn] {
				clientLocks[conn] = true
				readyClients++
				if readyClients == maxClients && lockReady {
					for clientConn := range clients {
						fmt.Fprintln(clientConn, "start")
					}
				}
			}
			clientLocksMu.Unlock()
		}

		if _, err := writer.WriteString(msg); err != nil {
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
	readyClients--
	delete(clients, conn)
	delete(clientLocks, conn)
	for clientConn := range clients {

		fmt.Fprintf(clientConn, "num_client#%v\n", clientCount)
	}
	mu.Unlock()
	fmt.Printf("Client %s disconnected\n", conn.RemoteAddr())
}
