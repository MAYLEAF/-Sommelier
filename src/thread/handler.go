package thread

import (
	"logger"
	"net"
	"sync"
)

type Handler struct {
	conn     *net.TCPConn
	value    []string
	Schedule sync.WaitGroup
	lock     sync.Mutex
	Send     chan string
	err      error
}

func (e *Handler) Create(serverAddr string, value []string) {
	server, _ := net.ResolveTCPAddr("tcp", serverAddr)
	e.conn, e.err = net.DialTCP("tcp", nil, server)
	e.value = value
	if e.err != nil {
		logger.Error("Fail to connect to Server err : %v", e.err)
	}
}

func (e *Handler) RequestMaker(actions map[string]interface{}) {
	threadContext := context{}
	Actions = actions
	threadContext.Initialize(e)
	threadContext.react(e)
	e.conn.Close()
}

func (e *Handler) Write(message []byte) {
	threadWriter := writer{}
	threadWriter.write(e, message)
}
