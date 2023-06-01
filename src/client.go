package main

import (
	"bufio"
	"fmt"
	"net"
)

func client(g *Game, ip string) error {
	port := "1234"

	var err error
	g.client_connection, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}

	fmt.Println("Client connecté")
	defer g.client_connection.Close()

	scanner := bufio.NewScanner(g.client_connection)
	for scanner.Scan() {
		token := scanner.Text()
		fmt.Println("receive:", token)
		if "stop" == token {
			fmt.Fprintln(g.client_connection, "stop")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func sendSpace(conn net.Conn) {
	fmt.Fprintln(conn, "space")
}