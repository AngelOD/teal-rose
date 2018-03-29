package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func SocketServer(port int) {
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Error starting server, %s", err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Println("Server running. Awaiting connections...")

	for {
		conn, err := listen.Accept()

		if err != nil {
			log.Fatalln(err)
			continue
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	var sb strings.Builder

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
	)

	log.Printf("Connection established (%s)\nAwaiting data...", conn.RemoteAddr())

ILOOP:

	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			sb.WriteString(data)
			if isTransportOver(data) {
				rd := ParseData(sb.String())

				log.Printf("Received: %+v\n", rd)

				if saveData {
					// TODO: Save data
					log.Println("ERROR! Unable to save data!")
				}

				sb.Reset()
			}
		default:
			log.Fatalf("Receive data failed: %s", err)
			return
		}
	}

	log.Printf("Done!")
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\n")
	return
}
