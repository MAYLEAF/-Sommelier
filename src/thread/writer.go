package thread

import (
	"io"
	"json"
	"logger"
	"os"
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
	_ = io.MultiWriter(thread.conn, os.Stdout)
	msg.SetEncoder(thread.conn)

	if err := msg.Encode(msg.Json()); nil != err {
		logger.Info("failed to write err: %v", err)
		return err
	}
	defer logger.Info("Request Message:" + string(content))

	return nil
}
