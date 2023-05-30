package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
	fmt.Fprintln(conn, "salut la team ")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		token := scanner.Text()
		fmt.Println("receive:", token)
		fmt.Fprintln(conn, "ack", token)
		if "stop" == token {
			fmt.Fprintln(conn, "stop")
			break
		}
	}
}
