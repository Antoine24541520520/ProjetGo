package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	dialAndListen()
}

func dialAndListen() {
	local := true

	ip := "172.21.65.21"

	if local {
		ip = "localhost"
	}

	port := "1234"
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	fmt.Fprintln(conn, "salut la team ")

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Fprintln(conn, "Ping")
			}
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		token := scanner.Text()
		fmt.Println("receive:", token)
		if "stop" == token {
			fmt.Fprintln(conn, "stop")
			break
		}
	}
}