package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func SocketServer(port int, prg *program) {
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		logger.Errorf("Error starting server, %s", err)
		return
	}

	defer listen.Close()

	logger.Info("Server running. Awaiting connections...")

	for {
		select {
		case <-prg.exit:
			return
		default:
			// Do nothing
		}

		listen.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Second))
		conn, err := listen.Accept()

		if err != nil {
			oerr, ok := err.(*net.OpError)

			if !ok || !oerr.Timeout() {
				logger.Error(err)
			}

			continue
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	var sb strings.Builder

	var (
		buf       = make([]byte, 1024)
		r         = bufio.NewReader(conn)
	)

	logger.Infof("Connection established (%s)", conn.RemoteAddr())
	logger.Info("Awaiting data...")

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

				if debugLog {
					log.Printf("Received: %+v\n", sb.String())
					log.Printf("Parsed: %+v\n\n", rd)
				}

				if saveData {
					if len(rd.Sensors) > 0 {
						rdStore <-rd
					}
				}

				sb.Reset()
			}
		default:
			if debugLog {
				log.Printf("ERROR! Receive data failed: %s\n", err)
			}

			return
		}
	}

	logger.Info("Done!")
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\n")
	return
}
