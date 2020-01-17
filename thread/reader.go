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

func (Reader *reader) read(thread *Handler) []byte {
	Reader.lock.Lock()
	defer Reader.lock.Unlock()

	msg := make(map[string]interface{})
	Reader.once.Do(func() {
		Reader.bufReader = bufio.NewReader(thread.conn)
		Reader.jsonDecoder = json.NewDecoder(Reader.bufReader)
	})

	if err := Reader.jsonDecoder.Decode(&msg); err != nil {
		logger.Error("Thread Reader error occur Err: %v User:%v", err, thread.value[0])
	}
	buf := jsonLib.Read(msg)

	logger.Info("Message from server - User=" + thread.value[0] + " " + string(buf) + "\n")
	return buf
}
