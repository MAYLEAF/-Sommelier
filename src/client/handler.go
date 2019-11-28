package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type handler struct {
	conn  net.Conn
	value []string
	err   error
}

func (e *handler) Create(serverAddr string, value []string) {
	e.conn, e.err = net.Dial("tcp", serverAddr)
	e.value = value
	if e.err != nil {
		log.Fatalf("Fail to connect to Server")
	}
}

func (e *handler) MakeRequest(Message string, ch chan string) {
	var result interface{}
	json.Unmarshal([]byte(Message), &result)
	result.(interface{}).(map[string]interface{})["uid"] = e.value[0]
	message, _ := json.Marshal(result)
	e.conn.Write([]byte(message))
	log.Print("Request" + string(message))
	ch <- e.value[0]
}

func (e *handler) ListenResponse(ch chan string) {

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
		buf = append(buf, tmp[:n]...)
	}

	fmt.Print("Message from server: " + string(buf) + "\n")

}
