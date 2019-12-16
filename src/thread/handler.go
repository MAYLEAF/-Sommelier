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

func (e *Handler) RequestMaker(actions map[string]interface{}) {
	log.Printf("Logger: requestMaker START handler=%v", e)
	defer log.Printf("Logger: requestMaker END handler=%v\n\n", e)
	threadContext := context{}
	threadContext.create(actions)
	threadContext.react(e)
}

func (e *Handler) Write(message string) {
	threadWriter := writer{}
	threadWriter.write(e, message)
}
