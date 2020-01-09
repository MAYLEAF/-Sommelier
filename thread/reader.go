package thread

import (
	"bufio"
	"github.com/MAYLEAF/Sommelier/json"
	"github.com/MAYLEAF/Sommelier/logger"
	"sync"
)

type reader struct {
	thread Handler
	lock   sync.Mutex
}

func (e *reader) read(thread *Handler) []byte {
	e.lock.Lock()
	defer e.lock.Unlock()

	msg := make(map[string]interface{})
	bufReader := bufio.NewReader(thread.conn)
	if err := json.Decode(bufReader, msg); err != nil {
		logger.Info("Thread Reader error occur Err: %v User:%v", err, thread.value[0])
	}
	buf := json.Read(msg)

	logger.Info("Message from server - User=" + thread.value[0] + " " + string(buf) + "\n")
	return buf
}
