package thread

import (
	"bytes"
	"json"
	"logger"
	"sync"
)

type writer struct {
	lock sync.Mutex
}

func (e *writer) write(thread *Handler, message []byte) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	defer logger.Info("Request Message:" + string(message))

	msg := make(map[string]interface{})
	byteReader := bytes.NewReader(message)
	json.Decode(byteReader, msg)

	msg["uid"] = thread.value[0]

	if err := json.Encode(thread.conn, msg); err != nil {
		logger.Info("Thread Encode error occur. Err: %v, Conn: %v", err, thread.conn)
		return err
	}
	return nil
}
