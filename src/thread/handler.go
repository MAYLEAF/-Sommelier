package thread

import (
	"logger"
	"net"
	"sync"
)

var logger = logger.Logger()

type Handler struct {
	conn     net.Conn
	value    []string
	Schedule sync.WaitGroup
	lock     sync.Mutex
	Send     chan string
	err      error
}

func (e *Handler) Create(serverAddr string, value []string) {
	e.conn, e.err = net.Dial("tcp", serverAddr)
	e.value = value
	if e.err != nil {
		logger.Error("Fail to connect to Server")
	}
}

func (e *Handler) RequestMaker(actions map[string]interface{}) {
	threadContext := context{}
	Actions = actions
	threadContext.Initialize()
	threadContext.react(e)
	e.conn.Close()
}

func (e *Handler) Write(message string) {
	threadWriter := writer{}
	threadWriter.write(e, message)
}
