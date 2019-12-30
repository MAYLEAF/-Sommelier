package thread

import (
	"logger"
	"net"
	"sync"
)

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

func (e *Handler) Write(message []byte) {
	threadWriter := writer{}
	threadWriter.write(e, message)
}
