package thread

import (
	"bufio"
	json "encoding/json"
	jsonLib "github.com/MAYLEAF/Sommelier/json"
	"github.com/MAYLEAF/Sommelier/logger"
	"sync"
)

type reader struct {
	thread      Handler
	jsonDecoder *json.Decoder
	bufReader   *bufio.Reader
	lock        sync.Mutex
	once        sync.Once
}

func (e *reader) read(thread *Handler) []byte {
	e.lock.Lock()
	defer e.lock.Unlock()

	msg := make(map[string]interface{})
	e.once.Do(func() {
		e.bufReader = bufio.NewReader(thread.conn)
		e.jsonDecoder = json.NewDecoder(e.bufReader)
	})

	if err := e.jsonDecoder.Decode(&msg); err != nil {
		logger.Info("Thread Reader error occur Err: %v User:%v", err, thread.value[0])
	}
	buf := jsonLib.Read(msg)

	logger.Info("Message from server - User=" + thread.value[0] + " " + string(buf) + "\n")
	return buf
}
