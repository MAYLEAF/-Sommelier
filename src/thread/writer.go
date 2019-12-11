package thread

import (
	"json"
	"log"
	"sync"
)

type writer struct {
	lock sync.Mutex
}

func (e *writer) write(thread *Handler, message string) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	msg := json.Json{}
	msg.Create(message)
	msg.Update("uid", thread.value[0])
	content, _ := msg.Read()

	if _, err := thread.conn.Write(content); nil != err {
		log.Printf("failed to write err: %v", err)
		return err
	}
	log.Print("Request Message:" + string(content))

	return nil
}
