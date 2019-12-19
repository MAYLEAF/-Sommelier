package thread

import (
	"json"
	"logger"
	"sync"
)

type writer struct {
	lock sync.Mutex
}

func (e *writer) write(thread *Handler, message string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	logger := logger.Logger()

	msg := json.Json{}
	msg.Create(message)
	msg.Update("uid", thread.value[0])
	content := msg.Read()

	if _, err := thread.conn.Write(content); nil != err {
		logger.Info("failed to write err: %v", err)
		return err
	}
	logger.Info("Request Message:" + string(content))

	return nil
}
