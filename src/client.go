package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func client(g *Game, ip string) error {
	port := "1234"

	var err error
	g.client_connection, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	defer g.client_connection.Close()

	// Send "Ping" message
	sendSpace(g.client_connection)

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
			break
		}
		if tokenSplited[0] == "num_client" {
			fmt.Println(tokenSplited[1])
			g.numClient = tokenSplited[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func sendSpace(conn net.Conn) {
	fmt.Fprintln(conn, "Ping")
}
