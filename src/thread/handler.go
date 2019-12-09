package thread

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Handler struct {
	conn     net.Conn
	value    []string
	Schedule sync.WaitGroup
	Send     chan string
	err      error
}

func (e *Handler) Create(serverAddr string, value []string) {
	e.conn, e.err = net.Dial("tcp", serverAddr)
	e.value = value
	if e.err != nil {
		log.Fatalf("Fail to connect to Server")
	}
}

func (e *Handler) RequestMaker() {
	log.Printf("Logger: handler.requestMaker() handler=%v", e)
	defer log.Printf("Logger: handler.requestMaker() handler=%v", e)
	threadReader := reader{}

	for {
		select {
		case msg := <-e.Send:
			ch := make(chan string, 2)
			time.Sleep(100 * time.Millisecond)
			if err := e.MakeRequest(msg, ch); nil != err {
				log.Printf("failed request err: %v", err)
			}
			threadReader.read(e)
			e.Schedule.Done()
			break
		}
	}
}

func (e *Handler) MakeRequest(Message string, ch chan string) error {

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
	log.Print("Request Message:" + string(message))

	ch <- e.value[0]
	return nil
}
