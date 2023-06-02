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
	maxClients = 4
)

var (
	clientCount    int
	readyClients   int
	lockReady      bool
	mu             sync.Mutex
	lockReadyMutex sync.Mutex
	clientLocksMu  sync.Mutex
	clients        = make(map[net.Conn]int)
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
		clients[conn] = clientCount
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

func broadcastMessage(sender net.Conn, message string) {
	mu.Lock()
	defer mu.Unlock()
	for clientConn := range clients {
		if clientConn != sender {
			fmt.Fprint(clientConn, message)
		}
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	clientID := clients[conn]

	fmt.Printf("Client %s connected\n", conn.RemoteAddr())

	id := fmt.Sprintf("id/%v\n", clientID)
	fmt.Fprintf(conn, id)

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

			value := splitMsg[1]

			spaceMsg := fmt.Sprintf("space/%v/%s\n", clientID, value)

			fmt.Printf("Client %v pressed space with value: %s \n", clientID, value)

			broadcastMessage(conn, spaceMsg)
		}

		if strings.HasPrefix(msg, "locked") {
			splitMsg := strings.Split(msg, "/") // ["locked", "1"]
			if len(splitMsg) < 2 {
				fmt.Println("Invalid message format")
				continue
			}
			value := splitMsg[1]
			lockedMsg := fmt.Sprintf("locked/%v/%v\n", clientID, value)

			broadcastMessage(conn, lockedMsg)

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
		}

		if strings.HasPrefix(msg, "finish") {
			splitMsg := strings.Split(strings.TrimSuffix(msg, "\n"), "/")
			if len(splitMsg) < 2 {
				fmt.Println("Invalid finish message format")
				continue
			}

			value := splitMsg[1]

			finishMsg := fmt.Sprintf("finish/%v/%s\n", clientID, value)

			fmt.Printf("Client %v finished with time: %s\n", clientID, value)

			broadcastMessage(conn, finishMsg)
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
