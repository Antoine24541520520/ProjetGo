package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
)

func client(ctx context.Context, g *Game, ip string) error {
    port := "1234"
    var d net.Dialer
    conn, err := d.DialContext(ctx, "tcp", ip+":"+port)
    if err != nil {
        return err
    }
    defer conn.Close()

    // Send "Ping" message
    sendSpace(conn)

    errChan := make(chan error, 1)
    go func() {
        scanner := bufio.NewScanner(conn)
        for scanner.Scan() {
            token := scanner.Text()
            fmt.Println("receive:", token)
            if "stop" == token {
                fmt.Fprintln(conn, "stop")
                break
            }
        }
        if err := scanner.Err(); err != nil {
            errChan <- err
        }
    }()

    select {
    case <-ctx.Done():
        return ctx.Err()
    case err := <-errChan:
        return err
    }
}

func sendSpace(conn net.Conn) {
	fmt.Fprintln(conn, "Ping")
}