package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func receiveMessage(g *Game) {
	scanner := bufio.NewScanner(g.client_connection)
	for scanner.Scan() {
		token := scanner.Text()
		tokenSplited := strings.Split(token, "#")
		fmt.Println("receive:", token)
		if token == "stop" {
			fmt.Fprintln(g.client_connection, "stop")
			break
		}
		if token == "ready" {
			g.lobbyReady = true
		}
		if token == "start" {
			g.start = true
		}
		if tokenSplited[0] == "num_client" {
			g.numClient = tokenSplited[1]
		}
		if strings.HasPrefix(token, "locked") {
			splitToken := strings.Split(token, "/")
			if len(splitToken) < 3 {
				fmt.Println("Invalid locked message format")
				continue
			}
			runnerID, err := strconv.Atoi(splitToken[1])
			if err != nil {
				fmt.Println("Invalid runner ID")
				continue
			}
			colorScheme, err := strconv.Atoi(splitToken[2])
			if err != nil {
				fmt.Println("Invalid colorScheme")
				continue
			}
			g.runners[runnerID-1].colorScheme = colorScheme
		}

		if strings.HasPrefix(token, "space") {
			splitToken := strings.Split(token, "/")
			if len(splitToken) < 3 {
				fmt.Println("Invalid space message format")
				continue
			}
			runnerID, err := strconv.Atoi(splitToken[1])
			if err != nil {
				fmt.Println("Invalid runner ID")
				continue
			}
			xpos, err := strconv.ParseFloat(splitToken[2], 64)
			if err != nil {
				fmt.Println("Invalid xpos")
				continue
			}
			g.runners[runnerID-1].xpos = xpos
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}

func sendSpace(conn net.Conn) {
	fmt.Fprintln(conn, "space/10")
	fmt.Println("space")
}

func sendLockChoice(conn net.Conn, colorScheme int) {
	msg := fmt.Sprintf("locked/%v", colorScheme)
	fmt.Fprintln(conn, msg)
}
