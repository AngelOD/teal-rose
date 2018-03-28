// +build go1.9

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Port int `cli:"port" usage:"Port number to host on" dft:"3333"`
	DbHost string `cli:"*dbhost" usage:"Host of DB to connect to" dft:"127.0.0.1"`
	DbName string `cli:"*dbname" usage:"Name of DB to use"`
	DbUser string `cli:"*dbuser" usage:"DB username" prompt:"Type DB username"`
	DbPass string `pw:"*dbpass" usage:"DB password" prompt:"Type DB password"`
}

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

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("%d\n%s\n%s\n%s\n%s\n", argv.Port, argv.DbHost, argv.DbName, argv.DbUser, argv.DbPass)
		return nil
	})
	//port := 3333
	//SocketServer(port)
}
