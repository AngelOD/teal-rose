// +build go1.9

package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "strconv"
    "github.com/teris-io/cli"
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

func main() {
    app := setupCli()

    os.Exit(app.Run(os.Args, os.Stdout))
    //port := 3333
    //SocketServer(port)
}

func setupCli() cli.App {
    optPort := cli.NewOption("port", "Port to host on").WithType(cli.TypeInt)
    optStoreData := cli.NewOption("save", "Store information in DB").WithChar('s').WithType(cli.TypeBool)

    app := cli.New("SW802F18 Test Server").WithOption(optPort).WithOption(optStoreData).WithAction(handleCli)

    return app
}

func handleCli(args []string, options map[string]string) int {
    fmt.Println("Args:", args)
    fmt.Println("Options:", options)

    return 0
}