// +build go1.9

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func SocketServer(port int) {
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	defer listen.Close()

	if err != nil {
		log.Fatalf("Socket listen port %d failed, %s", port, err)
		os.Exit(1)
	}

	log.Printf("Begin listen port: %d", port)

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

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
	)

ILOOP:

	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			log.Println("Receive: ", data)
		default:
			log.Fatalf("Receive data failed: %s", err)
			return
		}
	}

	log.Printf("Done!")
}
