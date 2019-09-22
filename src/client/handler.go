package client

import (
	"fmt"
	"io"
	"log"
	"net"
)

type handler struct {
	conn net.Conn
	err  error
}

func (e *handler) Create(serverAddr string, protocol string) {
	e.conn, e.err = net.Dial(protocol, serverAddr)
	if e.err != nil {
		log.Fatalf("Fail to connect to Server")
	}
}

func (e handler) MakeRequest(Message string) {
	e.conn.Write([]byte(Message))
	log.Print("Request" + Message)

	buf := make([]byte, 0, 16384)
	tmp := make([]byte, 256)

	for {
		n, err := e.conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				log.Print(err)
			}
			break
		}
		if n != 256 {
			buf = append(buf, tmp[:n]...)
			break
		}
		if n == 0 {
			break
		}
		buf = append(buf, tmp[:n]...)
	}

	fmt.Print("Message from server: " + string(buf) + "\n")

}


