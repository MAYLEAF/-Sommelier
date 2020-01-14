package thread

import (
	"github.com/MAYLEAF/Sommelier/logger"
	"net"
	"sync"
)

type Handler struct {
	conn     *net.TCPConn
	id       string
	value    []string
	Schedule sync.WaitGroup
	lock     sync.Mutex
	err      error
}

func (thread *Handler) Create(serverAddr string, value []string) {
	server, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		logger.Error("Fail to Handler Create Addr Resolve err : %v %v", err)
	}
	thread.conn, thread.err = net.DialTCP("tcp", nil, server)
	thread.value = value
	thread.id = value[0]
	if thread.err != nil {
		logger.Error("Fail to connect to Server err : %v %v %v", thread.err, server, serverAddr)
	}
}

func (thread *Handler) Attack(actions map[string]interface{}) {
	Context := context{}
	Actions = actions
	Context.Initialize(thread)
	Context.react(thread)
	thread.conn.Close()
}
