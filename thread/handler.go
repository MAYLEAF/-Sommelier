package thread

import (
	"github.com/MAYLEAF/Sommelier/logger"
	"net"
	"sync"
)

type Handler struct {
	conn     net.Conn
	id       string
	value    []string
	Schedule sync.WaitGroup
	lock     sync.Mutex
}

func New(rAddr string, value []string) *Handler {
	var err error
	newThread := Handler{}
	if newThread.conn, err = net.Dial("tcp", rAddr); err != nil {
		logger.Error("Fail to connect to Server err : %v", err)
	}
	newThread.value = value
	newThread.id = value[0]
	return &newThread
}

func (thread *Handler) Attack(actions map[string]interface{}) {
	Context := context{}
	Actions = actions
	Context.Initialize(thread)
	Context.react(thread)
	thread.conn.Close()
}
