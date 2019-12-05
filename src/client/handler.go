package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type handler struct {
	conn  net.Conn
	value []string
	send  chan string
	err   error
}

func (e *handler) Create(serverAddr string, value []string) {
	e.conn, e.err = net.Dial("tcp", serverAddr)
	e.value = value
	if e.err != nil {
		log.Fatalf("Fail to connect to Server")
	}
}

func (e *handler) test(messages []string, thread sync.WaitGroup) {
	go e.requestMaker()
	for _, message := range messages {
		ch := make(chan string, 1)
		e.MakeRequest(message, ch)
	}
	thread.Done()
}

func (e *handler) requestMaker() {
	for {
		select {
		case msg := <-e.send:
			ch := make(chan string, 2)
			if err := e.MakeRequest(msg, ch); nil != err {
				log.Printf("failed request err: %v", err)
			}
		}
	}
}

func (e *handler) MakeRequest(Message string, ch chan string) error {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	var result interface{}
	json.Unmarshal([]byte(Message), &result)
	result.(interface{}).(map[string]interface{})["uid"] = e.value[0]
	message, _ := json.Marshal(result)
	if _, err := e.conn.Write([]byte(message)); nil != err {
		return err
	}
	log.Print("Request" + string(message))
	ch <- e.value[0]
	return nil
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
