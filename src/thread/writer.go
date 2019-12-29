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

	msg := make(map[string]interface{})
	strReader := strings.NewReader(message)
	json.Decode(strReader, msg)

	msg["uid"] = thread.value[0]

	json.Encode(thread.conn, msg)
	defer logger.Info("Request Message:" + string(content))
	return nil
}
