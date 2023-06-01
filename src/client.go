package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func receiveMessage(g Game){
	scanner := bufio.NewScanner(g.client_connection)
	for scanner.Scan() {
		token := scanner.Text()
		tokenSplited := strings.Split(token, "#")
		fmt.Println("receive:", token)
		if "stop" == token {
			fmt.Fprintln(g.client_connection, "stop")
			break
		}
		if "ready" == token {
			g.lobbyReady = true
		}
		if "start" == token {
			g.start = true
		}
		if tokenSplited[0] == "num_client" {
			g.numClient = tokenSplited[1]
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}

func sendSpace(conn net.Conn) {
	fmt.Fprintln(conn, "space")
	fmt.Println("space")
}

func sendLockChoice(conn net.Conn) {
	fmt.Fprintln(conn, "locked")
}
