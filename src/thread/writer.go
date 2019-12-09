package thread

import (
	"encoding/json"
	"log"
	"sync"
)

type writer struct {
	lock sync.Mutex
}

func (e *writer) write(thread *Handler, Message string) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	var result interface{}
	json.Unmarshal([]byte(Message), &result)
	result.(interface{}).(map[string]interface{})["uid"] = thread.value[0]
	message, _ := json.Marshal(result)
	if _, err := thread.conn.Write([]byte(message)); nil != err {
		log.Printf("failed to write err: %v", err)
		return err
	}
	log.Print("Request Message:" + string(message))

	return nil
}
