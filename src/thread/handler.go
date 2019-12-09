package thread

import (
	"encoding/json"
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
	threadWriter := writer{}

	for {
		select {
		case msg := <-e.Send:
			time.Sleep(150 * time.Millisecond)
			threadWriter.write(e, msg)
			threadReader.read(e)
			e.Schedule.Done()
			break
		}
	}
}
