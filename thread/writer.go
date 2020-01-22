package thread

import (
	"bufio"
	"github.com/MAYLEAF/Sommelier/logger"
	"sync"
)

type writer struct {
	lock sync.Mutex
}

func (e *writer) write(thread *Handler, message []byte) error {
	thread.lock.Lock()
	defer thread.lock.Unlock()

	bufWriter := bufio.NewWriter(thread.conn)

	if _, err := bufWriter.Write(message); err != nil {
		logger.Error("Thread Writer Write error occur. Err: %v, Conn: %v", err, thread.value[0])
		return err
	}
	if err := bufWriter.Flush(); err != nil {
		logger.Error("Thread Writer Flush error occur. Err: %v, Conn: %v", err, thread.value[0])
		return err
	}
	return nil
}
