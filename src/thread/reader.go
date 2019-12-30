package thread

import (
	"json"
	"logger"
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
	json.Decode(thread.conn, msg)
	buf := json.Read(msg)

	logger.Info("Message from server - User=" + thread.value[0] + " " + string(buf) + "\n")
	return buf
}
