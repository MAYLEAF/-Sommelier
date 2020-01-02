package thread

import (
	_ "bufio"
	"json"
	"logger"
	"sync"
)

type reader struct {
	thread Handler
	lock   sync.Mutex
}

func (e *reader) read(thread *Handler) []byte {
	thread.lock.Lock()
	defer thread.lock.Unlock()

	msg := make(map[string]interface{})
	if err := json.Decode(thread.conn, msg); err != nil {
		logger.Info("Thread Reader error occur Err: %v User:%v", err, thread.value[0])
	}
	buf := json.Read(msg)

	logger.Info("Message from server - User=" + thread.value[0] + " " + string(buf) + "\n")
	return buf
}
