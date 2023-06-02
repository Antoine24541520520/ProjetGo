package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func receiveMessage(g *Game) {
	scanner := bufio.NewScanner(g.client_connection)
	for scanner.Scan() {
		token := strings.TrimSuffix(scanner.Text(), "\n")
		tokenSplited := strings.Split(token, "#")
		fmt.Println("receive:", token)
		if token == "stop" {
			fmt.Fprint(g.client_connection, "stop")
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
		if strings.HasPrefix(token, "id") {
			splitToken := strings.Split(token, "/")
			if len(splitToken) < 2 {
				fmt.Println("Invalid locked message format")
				continue
			}
			posRunner, err := strconv.Atoi(splitToken[1])
			if err != nil {
				fmt.Println("Invalid runner ID")
				continue
			}
			g.posMainRunner = posRunner - 1
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

		if strings.HasPrefix(token, "finish") {
			splitToken := strings.Split(strings.TrimSuffix(token, "\n"), "/")
			if len(splitToken) < 2 {
				fmt.Println("Invalid finish message format")
				continue
			}
			runnerID, err := strconv.Atoi(splitToken[1])
			if err != nil {
				fmt.Println("Invalid runner ID")
				continue
			}
			finishTime, err := time.ParseDuration(splitToken[2] + "ms")
			if err != nil {
				fmt.Println("Invalid finish time")
				continue
			}
			g.runners[runnerID-1].runTime = finishTime
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}

func sendFinishTime(conn net.Conn, finishTime time.Duration) {
	msg := fmt.Sprintf("finish/%d\n", finishTime.Milliseconds())
	fmt.Fprint(conn, msg)
}

func sendSpace(conn net.Conn, newPos float64) {
	msg := fmt.Sprintf("space/%f\n", newPos)
	fmt.Fprint(conn, msg)
}

func sendLockChoice(conn net.Conn, colorScheme int) {
	msg := fmt.Sprintf("locked/%v\n", colorScheme)
	fmt.Fprint(conn, msg)
}
