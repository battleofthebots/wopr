package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:4000")
	if err != nil {
		log.Fatal(err)
	}
	j := 0
	for {
		// set SetReadDeadline
		err := conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		if err != nil {
			log.Fatal("SetReadDeadline failed:", err)
		}
		input, err := io.ReadAll(conn)
		i := string(input)
		fmt.Print(i)

		if strings.Contains(i, "bested") {
			fmt.Println("[BOT] I win!! Forking over a shell")
			break
		}
		if strings.Contains(i, "<< enter move [0-8]:") {
			fmt.Println("\n[BOT] its my turn!", j)
			fmt.Fprintln(conn, j)
			j = (j + 1) % 9

		}
		if strings.Contains(i, "[WOPR] let's play again") {
			j = 0
			continue
		}
	}

	go func() {
		for {
			conn.SetReadDeadline(time.Now().Add(time.Second))
			io.Copy(os.Stdout, conn)
		}
	}()
	for {
		io.Copy(conn, os.Stdin)
	}
}
