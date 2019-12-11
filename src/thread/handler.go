package thread

import (
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
	defer log.Printf("Logger: handler.requestMaker() handler=%v\n\n", e)
	threadReader := reader{}
	threadWriter := writer{}

	for {
		select {
		case msg := <-e.Send:
			time.Sleep(150 * time.Millisecond)
			threadWriter.write(e, msg)
			threadReader.read(e, msg)
			e.Schedule.Done()
			break
		}
	}
}

func (e *Handler) Write(message string) {
	threadWriter := writer{}
	threadWriter.write(e, message)
}
