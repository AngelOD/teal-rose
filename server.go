package main

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
	"log"
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
		saveQueue = make([]RadioData, 0, saveEvery)
	)

	logger.Infof("Connection established (%s)\nAwaiting data...", conn.RemoteAddr())

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
					if len(rd.Sensors) > 0 {
						saveQueue = append(saveQueue, rd)

						if len(saveQueue) >= saveEvery {
							if StoreData(saveQueue) {
								log.Println("Data saved successfully!")
								saveQueue = nil
							} else {
								log.Println("ERROR! Unable to save data!")
							}
						}
					}
				}

				sb.Reset()
			}
		default:
			log.Printf("ERROR! Receive data failed: %s\n", err)
			return
		}
	}

	logger.Info("Done!")
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\n")
	return
}
