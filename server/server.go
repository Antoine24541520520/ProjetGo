package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

const (
	maxClients = 2
)

var (
	clientCount    int
	readyClients   int
	lockReady      bool
	mu             sync.Mutex
	lockReadyMutex sync.Mutex
	clientLocksMu  sync.Mutex
	clients        = make(map[net.Conn]int) // Changes here
	clientLocks    = make(map[net.Conn]bool)
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
		clients[conn] = clientCount // Changes here
		clientLocksMu.Lock()
		clientLocks[conn] = false
		clientLocksMu.Unlock()
		fmt.Printf("num_client#%v", clientCount)
		for clientConn, id := range clients { // Changes here

			fmt.Fprintf(clientConn, "num_client#%v\n", id) // Changes here
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
	defer wg.Done()

	fmt.Printf("Client %s connected\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading message from %s: %v\n", conn.RemoteAddr(), err)
			}
			break
		}

		if strings.HasPrefix(msg, "space") {
			splitMsg := strings.Split(strings.TrimSuffix(msg, "\n"), "/")
			if len(splitMsg) < 2 {
				fmt.Println("Invalid space message format")
				continue
			}

			clientID := clients[conn]
			value := splitMsg[1]

			spaceMsg := fmt.Sprintf("space/%v/%s\n", clientID, value)

			fmt.Printf("Client %v pressed space with value: %s \n", clientID, value)

			mu.Lock()
			for clientConn := range clients {
				if clientConn != conn {
					fmt.Fprint(clientConn, spaceMsg)
				}
			}
			mu.Unlock()
		}

		if strings.HasPrefix(msg, "locked") {
			splitMsg := strings.Split(msg, "/") // ["locked", "1"]
			if len(splitMsg) < 2 {
				fmt.Println("Invalid message format")
				continue
			}
			value := splitMsg[1]
			clientID := clients[conn]
			lockedMsg := fmt.Sprintf("locked/%v/%v\n", clientID, value)

			clientLocksMu.Lock()
			if !clientLocks[conn] {
				clientLocks[conn] = true
				readyClients++
				if readyClients == maxClients && lockReady {
					for clientConn := range clients {
						fmt.Fprintln(clientConn, "start")
					}
				}
				clientLocksMu.Unlock()
			}
			mu.Lock()
			for clientConn := range clients {
				if clientConn != conn {
					fmt.Fprintf(clientConn, lockedMsg)
				}
			}
			mu.Unlock()

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
	conn.Close()
}
